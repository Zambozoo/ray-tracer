package tracer

import (
	"sync"
)

func doParallel[T any](numWorkers int, queue chan T, wg *sync.WaitGroup, f func(T)) {
	for i := 0; i < numWorkers; i++ {
		go func() {
			for item := range queue {
				f(item)
			}
		}()
	}

	if wg != nil {
		wg.Wait()
	}
}
