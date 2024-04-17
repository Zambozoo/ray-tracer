package tracer

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"v0/io"
	"v0/primitive"
)

type WaitGroup struct {
	wg      sync.WaitGroup
	counter atomic.Int64
}

func (wg *WaitGroup) Add(i int) {
	wg.wg.Add(i)
	wg.counter.Add(int64(i))
}

func (wg *WaitGroup) Done() {
	wg.wg.Done()
	wg.counter.Add(-1)
}

func (wg *WaitGroup) Load() int64 {
	return wg.counter.Load()
}

func RayTrace(input *io.Input) {
	fmt.Print("--Beginning RayTracing--\n")
	pf := io.NewPPM(io.Env.OutputFilePath,
		primitive.Vector2I{
			X: input.ImageWidth,
			Y: int(float64(input.ImageWidth) / input.AspectRatio),
		},
	)
	pf.InitFile()

	for row := 0; row < pf.Height(); row++ {
		rayQueue, rayWaitGroup := initializeRayQueue(pf, input, row)
		pixelQueue, pixelWaitGroup := initializePixelQueue(pf, input)
		t := time.NewTicker(30 * time.Second)
		go func() {
			for range t.C {
				fmt.Printf("\t%d Remaining rays for row %d\n", rayWaitGroup.Load(), row)
			}
		}()

		doParallel(input.NumWorkers, rayQueue, &rayWaitGroup.wg,
			cast(input, pixelQueue, pixelWaitGroup, rayWaitGroup),
		)

		t.Stop()

		fmt.Printf("--Finished RayCasting for row %d--\n\n", row)

		close(rayQueue)
		pixelWaitGroup.Wait()
		close(pixelQueue)

		fmt.Printf("--Finished Summing Colors for row %d--\n\n", row)

		t.Stop()
		pf.Write(row)
	}

	fmt.Print("--Finished Writing to File--\n")
}

type weightedRay struct {
	weight      float64
	ray         primitive.Ray3
	depth       int
	attenuation primitive.Color
	pixel       *primitive.AtomicPixel
}

type coloredPixel struct {
	pixel *primitive.AtomicPixel
	color primitive.Color
}

func initializePixelQueue(pf *io.PPMFile, input *io.Input) (chan coloredPixel, *sync.WaitGroup) {
	bufferSize := 2 * pf.Width() * input.MaxDepth * input.AntiAliasCount
	coloredPixels := make(chan coloredPixel, bufferSize)
	pixelQueueWaitGroup := &sync.WaitGroup{}

	go func() {
		doParallel(input.NumWorkers, coloredPixels, nil, func(cp coloredPixel) {
			cp.pixel.Add(cp.color)
			pixelQueueWaitGroup.Done()
		})
	}()

	return coloredPixels, pixelQueueWaitGroup
}

func initializeRayQueue(pf *io.PPMFile, input *io.Input, row int) (chan weightedRay, *WaitGroup) {
	bufferSize := pf.Width() * input.MaxDepth * input.AntiAliasCount * input.AntiAliasCount
	weightedRays := make(chan weightedRay, bufferSize)

	go func() {
		camera := input.Camera
		viewPlane := camera.ViewPlaneRect(input.AspectRatio)
		dRay := input.Camera.ViewPlaneWidth / float64(input.ImageWidth)

		normX := camera.NormX()
		normY := camera.NormY().Scale(-1)

		zeroVector := viewPlane.TopLeft.Add(normX.Add(normY).Scale(0.5 * dRay))
		ddRay := dRay / float64(input.AntiAliasCount)
		for _, pixelOffset := range pf.PixelIndicies(row) {
			ap := pf.GetPixel(pixelOffset.X, pixelOffset.Y)
			for i := 0; i < input.AntiAliasCount; i++ {
				for j := 0; j < input.AntiAliasCount; j++ {
					xScale := dRay*(float64(pixelOffset.X)) + ddRay*(float64(i)+rand.Float64())
					yScale := dRay*(float64(pixelOffset.Y)) + ddRay*(float64(j)+rand.Float64())

					ray := primitive.NewRay3(camera.LookRay.Src,
						zeroVector.
							Add(normX.Scale(xScale)).
							Add(normY.Scale(yScale)).Subtract(camera.LookRay.Src).Norm(),
					)
					weightedRays <- weightedRay{
						weight:      1.0 / float64(input.AntiAliasCount*input.AntiAliasCount),
						ray:         ray,
						depth:       1,
						pixel:       ap,
						attenuation: primitive.White,
					}
				}
			}
		}
	}()

	rayWaitGroup := &WaitGroup{}
	rayWaitGroup.Add(pf.Width() * input.AntiAliasCount * input.AntiAliasCount)

	return weightedRays, rayWaitGroup
}
