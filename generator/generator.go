package generator

import "errors"

var Exhausted = errors.New("generator exhausted")

func FromSlice(theSlice []int) func() (int, error) {
	currentIndex := 0

	return func() (int, error) {
		if currentIndex >= len(theSlice) {
			return 0, Exhausted
		}

		value := theSlice[currentIndex]
		currentIndex++

		return value, nil
	}
}
