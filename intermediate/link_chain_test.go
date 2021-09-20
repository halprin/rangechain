package intermediate

import (
	"github.com/halprin/slice-chain/generator"
	"github.com/halprin/slice-chain/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter(t *testing.T) {
	inputSlice    := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
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
