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
	input := [...]string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	chain := FromArray(input)

	assert.ElementsMatch(t, chain.Slice(), input)
}

func TestFromChannel(t *testing.T) {
	innerInput := []string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	input := createTestChannel(innerInput)
	chain := FromChannel(input)

	assert.ElementsMatch(t, chain.Slice(), innerInput)
}

func createTestChannel(stringSlice []string) chan interface{} {
	stringChannel := make(chan interface{})

	go func() {
		for _, currentString := range stringSlice {
			stringChannel <- currentString
		}
		close(stringChannel)
	}()

	return stringChannel
}
