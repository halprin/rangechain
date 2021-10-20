package intermediate

import (
	"github.com/halprin/rangechain/generator"
	"github.com/halprin/rangechain/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapParallel(t *testing.T) {
	inputSlice := []string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	var expectedOutput []interface{}
	for _, stringValue := range inputSlice {
		expectedOutput = append(expectedOutput, len(stringValue))
	}

	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	mapFunction := func(value interface{}) interface{} {
		stringValue := value.(string)
		return len(stringValue)
	}

	actualSlice := link.MapParallel(mapFunction).Slice()

	assert.Equal(t, expectedOutput, actualSlice)
}

func TestFilterParallel(t *testing.T) {
	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	expectedSlice := helper.InterfaceSlice([]int{7, 9, 6, 8})
	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	filterFunction := func(value interface{}) bool {
		intValue := value.(int)
		return intValue > 5
	}

	actualSlice := link.FilterParallel(filterFunction).Slice()

	assert.Equal(t, expectedSlice, actualSlice)
}
