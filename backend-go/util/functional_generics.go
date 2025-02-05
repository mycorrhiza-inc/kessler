package util

import (
	"fmt"
	"sync"
)

func Map[T any, R any](input_list []T, f func(T) R) []R {
	results := make([]R, len(input_list))
	for i, item := range input_list {
		results[i] = f(item)
	}
	return results
}

func MapErrorDiscard[T any, R any](input_list []T, f func(T) (R, error)) []R {
	results := make([]R, len(input_list))
	for i, item := range input_list {
		result, err := f(item)
		if err != nil {
			fmt.Println("Error in Generic Map, told to discard: ", err)
		}
		results[i] = result
	}
	return results
}

func ConcurrentMapError[T any, R any](input_list []T, f func(T) (R, error), workers int) ([]R, error) {
	if workers == 0 {
		workers = len(input_list)
	}
	sem := make(chan struct{}, workers) // Semaphore to limit concurrency to 3
	errChan := make(chan error, len(input_list))
	resultChan := make(chan R, len(input_list))
	var wg sync.WaitGroup

	for index, value := range input_list {

		sem <- struct{}{} // Acquire semaphore slot
		wg.Add(1)

		go func(imbed_value T) {
			defer func() {
				<-sem // Release semaphore slot
				wg.Done()
			}()
			result, err := f(imbed_value)
			resultChan <- result
			if err != nil {
				errChan <- err
			}
		}(value)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Collect the first error encountered, if any
	var firstErr error
	close(errChan) // Close the channel to safely iterate
	for err := range errChan {
		if firstErr == nil {
			firstErr = err
		}
	}
	close(resultChan) // Close the channel to safely iterate
	results := make([]R, 0, len(input_list))
	for result := range resultChan {
		results = append(results, result)
	}

	return results, firstErr
}
