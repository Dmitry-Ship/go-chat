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

func MergeChannels[T any](ctx context.Context, channels ...chan T) <-chan T {
	var wg sync.WaitGroup
	wg.Add(len(channels))
	outChannel := make(chan T)

	worker := func(channel <-chan T) {
		defer wg.Done()
		for item := range channel {
			select {
			case outChannel <- item:
			case <-ctx.Done():
				return
			}
		}

	}

	for _, channel := range channels {
		go worker(channel)
	}

	go func() {
		wg.Wait()
		close(outChannel)
	}()

	return outChannel
}
