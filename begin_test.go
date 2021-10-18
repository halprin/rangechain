package slice_chain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromSlice(t *testing.T) {
	inputSlice := []string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	chain := FromSlice(inputSlice)

	assert.ElementsMatch(t, chain.Slice(), inputSlice)
}

func TestFromArray(t *testing.T) {
	input := [9]string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	chain := FromArray(input)

	assert.ElementsMatch(t, chain.Slice(), input)
}
