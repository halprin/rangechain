package intermediate

import (
	"github.com/halprin/slice-chain/generator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter(t *testing.T) {
	inputSlice    := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	expectedSlice := []int{7, 9, 6, 8}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	filterFunction := func(value int) bool {
		return value > 5
	}

	actualSlice := link.filter(filterFunction).slice()

	assert.Equal(t, expectedSlice, actualSlice)
}
