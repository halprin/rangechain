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
