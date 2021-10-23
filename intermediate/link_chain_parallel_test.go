package intermediate

import (
	"errors"
	"github.com/halprin/rangechain/generator"
	"github.com/halprin/rangechain/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapParallel(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	var expectedOutput []interface{}
	for _, stringValue := range inputSlice {
		expectedOutput = append(expectedOutput, len(stringValue))
	}

	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	mapFunction := func(value interface{}) (interface{}, error) {
		stringValue := value.(string)
		return len(stringValue), nil
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
	link := NewLink(generation)

	mapFunction := func(value interface{}) (interface{}, error) {
		stringValue := value.(string)
		if stringValue == errorValue {
			return 0, expectedError
		}
		return len(stringValue), nil
	}

	_, err := link.MapParallel(mapFunction).Slice()

	assert.Equal(expectedError, err)
}

func TestFilterParallel(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	expectedSlice := helper.InterfaceSlice([]int{7, 9, 6, 8})
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	filterFunction := func(value interface{}) bool {
		intValue := value.(int)
		return intValue > 5
	}

	actualSlice, err := link.FilterParallel(filterFunction).Slice()

	assert.Equal(expectedSlice, actualSlice)
	assert.Nil(err)
}
