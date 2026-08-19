package rangechain

import (
	"errors"
	"testing"

	"github.com/halprin/rangechain/internal/generator"
	"github.com/stretchr/testify/assert"
)

func TestMapParallel(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	var expectedOutput []int
	for _, stringValue := range inputSlice {
		expectedOutput = append(expectedOutput, len(stringValue))
	}

	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	mapFunction := func(value string) (int, error) {
		return len(value), nil
	}

	actualSlice, err := link.MapParallel(mapFunction).Slice()

	assert.Equal(expectedOutput, actualSlice)
	assert.Nil(err)
}

func TestMapParallelHasError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("an example error")
	errorValue := "Do"
	inputSlice := []string{"DogCows", "goes", "Moof!", errorValue, "you", "like", "Clarus", "the", "DogCow?"}

	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	mapFunction := func(value string) (int, error) {
		if value == errorValue {
			return 0, expectedError
		}
		return len(value), nil
	}

	_, err := link.MapParallel(mapFunction).Slice()

	assert.Equal(expectedError, err)
}

func TestFilterParallel(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	expectedSlice := []int{7, 9, 6, 8}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	filterFunction := func(value int) (bool, error) {
		return value > 5, nil
	}

	actualSlice, err := link.FilterParallel(filterFunction).Slice()

	assert.Equal(expectedSlice, actualSlice)
	assert.Nil(err)
}

func TestFilterParallelHasError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("an example error")
	errorValue := 9
	inputSlice := []int{7, 4, 2, 3, errorValue, 5, 6, 0, 8, 1}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	filterFunction := func(value int) (bool, error) {
		if value == errorValue {
			return false, expectedError
		}
		return value > 5, nil
	}

	_, err := link.FilterParallel(filterFunction).Slice()

	assert.Equal(expectedError, err)
}
