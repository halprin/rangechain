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

func TestFlattenWithSliceAndMap(t *testing.T) {
	innerMap := map[int]int{
		4: 5,
		6: 7,
	}
	inputSlice := []interface{}{[]int{1, 2, 3}, innerMap, []int{7, 8, 9}}
	//notice that the expected just puts the map right back into the slice without expanding it
	//in the future we'll support expanding other containers
	expectedSlice := []interface{}{1, 2, 3, innerMap, 7, 8, 9}

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
