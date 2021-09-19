package intermediate


import (
	"github.com/halprin/slice-chain/generator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlice(t *testing.T) {
	expectedSlice := []int{987, 8, 26}
	generation := generator.FromSlice(expectedSlice)
	link := NewLink(generation)

	actualSlice := link.Slice()

	assert.Equal(t, expectedSlice, actualSlice)
}

func TestForEach(t *testing.T) {
	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	var seenItems []int
	forEachFunction := func(value int) {
		seenItems = append(seenItems, value)
	}
	link.ForEach(forEachFunction)

	assert.ElementsMatch(t, inputSlice, seenItems)
}

func TestCount(t *testing.T) {
	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	actualCount := link.Count()

	assert.Equal(t, len(inputSlice), actualCount)
}

func TestFirst(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	actualFirst := link.First()

	assert.NotNil(actualFirst)
	assert.Equal(inputSlice[0], *actualFirst)
}

func TestFirstWithEmptySlice(t *testing.T) {
	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	actualFirst := link.First()

	assert.Nil(t, actualFirst)
}

func TestLast(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	actualLast := link.Last()

	assert.NotNil(actualLast)
	assert.Equal(inputSlice[len(inputSlice) - 1], *actualLast)
}

func TestLastWithEmptySlice(t *testing.T) {
	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	actualLast := link.Last()

	assert.Nil(t, actualLast)
}
