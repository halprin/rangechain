package intermediate

import (
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

	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	mapFunction := func(value interface{}) interface{} {
		stringValue := value.(string)
		return len(stringValue)
	}

	actualSlice, err := link.MapParallel(mapFunction).Slice()

	assert.Equal(expectedOutput, actualSlice)
	assert.Nil(err)
}

func TestFilterParallel(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	expectedSlice := helper.InterfaceSlice([]int{7, 9, 6, 8})
	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	filterFunction := func(value interface{}) bool {
		intValue := value.(int)
		return intValue > 5
	}

	actualSlice, err := link.FilterParallel(filterFunction).Slice()

	assert.Equal(expectedSlice, actualSlice)
	assert.Nil(err)
}
