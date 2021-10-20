package intermediate

import (
	"github.com/halprin/slice-chain/generator"
	"github.com/halprin/slice-chain/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap(t *testing.T) {
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

	actualSlice := link.Map(mapFunction).Slice()

	assert.Equal(t, expectedOutput, actualSlice)
}

func TestFilter(t *testing.T) {
	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	expectedSlice := helper.InterfaceSlice([]int{7, 9, 6, 8})
	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	filterFunction := func(value interface{}) bool {
		intValue := value.(int)
		return intValue > 5
	}

	actualSlice := link.Filter(filterFunction).Slice()

	assert.Equal(t, expectedSlice, actualSlice)
}

func TestSkip(t *testing.T) {
	inputSlice    := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	howManyToSkip := 3
	actualSlice := link.Skip(howManyToSkip).Slice()

	assert.Equal(t, helper.InterfaceSlice(inputSlice[howManyToSkip:]), actualSlice)
}

func TestSkipLargerThanSlice(t *testing.T) {
	inputSlice    := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	actualSlice := link.Skip(len(inputSlice) + 1).Slice()

	assert.Equal(t, []interface{}{}, actualSlice)
}

func TestLimit(t *testing.T) {
	inputSlice    := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	howManyToKeep := 6
	actualSlice := link.Limit(howManyToKeep).Slice()

	assert.Equal(t, helper.InterfaceSlice(inputSlice[:howManyToKeep]), actualSlice)
}

func TestLimitLargerThanSlice(t *testing.T) {
	inputSlice    := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	actualSlice := link.Limit(len(inputSlice) + 1).Slice()

	assert.Equal(t, helper.InterfaceSlice(inputSlice), actualSlice)
}

func TestDistinct(t *testing.T) {
	inputSlice := []int{7, 4, 2, 7, 3, 7, 9, 5, 5, 2, 6, 0, 8, 1}
	expectedSlice := helper.InterfaceSlice([]int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1})
	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	actualSlice := link.Distinct().Slice()

	assert.Equal(t, expectedSlice, actualSlice)
}

func TestFlattenWithSliceOfSlice(t *testing.T) {
	inputSlice := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	expectedSlice := helper.InterfaceSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	actualSlice := link.Flatten().Slice()

	assert.Equal(t, expectedSlice, actualSlice)
}

func TestFlattenWithSliceMix(t *testing.T) {
	inputSlice := []interface{}{[]int{1, 2, 3}, 4, []int{7, 8, 9}}
	expectedSlice := helper.InterfaceSlice([]int{1, 2, 3, 4, 7, 8, 9})
	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	actualSlice := link.Flatten().Slice()

	assert.Equal(t, expectedSlice, actualSlice)
}

func TestFlattenWithArray(t *testing.T) {
	inputSlice := []interface{}{[...]int{1, 2, 3}, 4, [...]int{7, 8, 9}}
	expectedSlice := helper.InterfaceSlice([]int{1, 2, 3, 4, 7, 8, 9})
	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	actualSlice := link.Flatten().Slice()

	assert.Equal(t, expectedSlice, actualSlice)
}

func TestFlattenWithChannel(t *testing.T) {
	firstChannel := createTestChannel([]int{1, 2, 3})
	secondChannel := createTestChannel([]int{7, 8, 9})

	inputSlice := []interface{}{firstChannel, 4, secondChannel}
	expectedSlice := helper.InterfaceSlice([]int{1, 2, 3, 4, 7, 8, 9})
	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	actualSlice := link.Flatten().Slice()

	assert.Equal(t, expectedSlice, actualSlice)
}

func TestFlattenWithSliceAndMap(t *testing.T) {
	key1 := 4
	value1 := 5
	key2 := 6
	value2 := 7

	innerMap := map[int]int{
		key1: value1,
		key2: value2,
	}
	inputSlice := []interface{}{[]int{1, 2, 3}, innerMap, []int{7, 8, 9}}
	expectedSlice := []interface{}{
		1,
		2,
		3,
		generator.MapTuple{
			Key:   key1,
			Value: value1,
		},
		generator.MapTuple{
			Key:   key2,
			Value: value2,
		},
		7,
		8,
		9,
	}

	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	actualSlice := link.Flatten().Slice()

	assert.Equal(t, expectedSlice, actualSlice)
}

func TestSort(t *testing.T) {
	inputSlice    := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	expectedSlice := helper.InterfaceSlice([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	returnLessFunction := func(sliceToSort []interface{}) func(int, int) bool {
		return func(i int, j int) bool {
			iItem := sliceToSort[i].(int)
			jItem := sliceToSort[j].(int)
			return iItem < jItem
		}
	}
	actualSlice := link.Sort(returnLessFunction).Slice()

	assert.Equal(t, expectedSlice, actualSlice)
}

func TestReverse(t *testing.T) {
	inputSlice    := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	expectedSlice := helper.InterfaceSlice([]int{1, 8, 0, 6, 5, 9, 3, 2, 4, 7})
	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	actualSlice := link.Reverse().Slice()

	assert.Equal(t, expectedSlice, actualSlice)
}

func createTestChannel(intSlice []int) chan interface{} {
	intChannel := make(chan interface{})

	go func() {
		for _, currentInt := range intSlice {
			intChannel <- currentInt
		}
		close(intChannel)
	}()

	return intChannel
}
