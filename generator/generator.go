package generator

import "errors"

var Exhausted = errors.New("generator exhausted")

func FromSlice(slice []interface{}) func() (interface{}, error) {
	currentIndex := 0

	return func() (interface{}, error) {
		if currentIndex >= len(slice) {
			return 0, Exhausted
		}

		value := slice[currentIndex]
		currentIndex++

		return value, nil
	}
}
