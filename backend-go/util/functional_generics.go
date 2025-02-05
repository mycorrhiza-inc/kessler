package util

import "fmt"

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
