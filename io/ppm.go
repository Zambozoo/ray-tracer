package io

import (
	"fmt"
	"math"
	"os"
	"time"
	"v0/primitive"

	"go.uber.org/zap"
)

type PPMFile struct {
	filePath string
	pixels   [][]primitive.AtomicPixel
}

func NewPPM(filePath string, rect primitive.Vector2I) *PPMFile {
	if rect.X == 0 || rect.Y == 0 {
		Logger.Panic("invalid ppm size",
			zap.Int("width", rect.X),
			zap.Int("height", rect.Y),
		)
	}
	pf := &PPMFile{
		filePath: filePath,
		pixels:   make([][]primitive.AtomicPixel, rect.X),
	}

	for i := range pf.pixels {
		pf.pixels[i] = make([]primitive.AtomicPixel, rect.Y)
	}

	return pf
}

func (pf *PPMFile) Width() int {
	return len(pf.pixels)
}

func (pf *PPMFile) Height() int {
	return len(pf.pixels[0])
}

func (pf *PPMFile) SetPixel(x, y int, p primitive.Color) {
	pf.pixels[x][y].Pixel = p
}

func (pf *PPMFile) GetPixel(x, y int) *primitive.AtomicPixel {
	return &(pf.pixels[x][y])
}

var f *os.File

func (pf *PPMFile) InitFile() {
	fileErrFunc := func(msg string, err error) {
		Logger.Panic(msg,
			zap.String("filePath", pf.filePath),
			zap.Error(err),
		)
	}

	var err error
	f, err = os.Create(pf.filePath)
	if err != nil {
		fileErrFunc("file creation err", err)
	} else if err := f.Truncate(0); err != nil {
		fileErrFunc("file truncate err", err)
	}

	width := len(pf.pixels)
	height := len(pf.pixels[0])

	if _, err := f.WriteString(fmt.Sprintf("P3\n%d\t%d\n255\n", width, height)); err != nil {
		fileErrFunc("file write err", err)
	}
}
func (pf *PPMFile) Write(row int) {
	fileErrFunc := func(msg string, err error) {
		Logger.Panic(msg,
			zap.String("filePath", pf.filePath),
			zap.Error(err),
		)
	}
	width := len(pf.pixels)

	i := 0
	t := time.NewTicker(time.Second)
	go func() {
		for range t.C {
			fmt.Printf("\t%d/%d Remaining pixels\n", i, width)
		}
	}()

	rowString := ""
	for x := 0; x < width; x++ {
		i++
		p := pf.pixels[x][row].Pixel
		pixelF := primitive.Color{
			R: math.Sqrt(p.R),
			G: math.Sqrt(p.G),
			B: math.Sqrt(p.B),
		}
		pixel := pixelF.Truncate()
		spacer := "\t"
		if x == width-1 {
			spacer = "\n"
		}
		rowString += fmt.Sprintf("%d %d %d%s", pixel.R, pixel.G, pixel.B, spacer)
	}
	if _, err := f.WriteString(rowString); err != nil {
		fileErrFunc("file write err", err)
	}

	t.Stop()
}

func (pf *PPMFile) PixelIndicies(row int) []primitive.Vector2I {
	indicies := make([]primitive.Vector2I, 0, pf.Width())
	for x := 0; x < pf.Width(); x++ {
		indicies = append(indicies, primitive.Vector2I{X: x, Y: row})
	}

	return indicies
}
