package infra

import (
	"context"
	"sync"
)

func FanOut[T any](numberOfWorkers int, worker func() (chan T, chan error)) ([]chan T, []chan error) {
	var outputChanSlice []chan T
	var errorChanSlice []chan error

	for i := 0; i < numberOfWorkers; i++ {
		outputChan, errorChan := worker()
		outputChanSlice = append(outputChanSlice, outputChan)
		errorChanSlice = append(errorChanSlice, errorChan)
	}

	return outputChanSlice, errorChanSlice
}

func MergeChannels[T any](ctx context.Context, cs ...chan T) <-chan T {
	var wg sync.WaitGroup
	wg.Add(len(cs))
	out := make(chan T)

	output := func(c <-chan T) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}

	}

	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
