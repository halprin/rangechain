package rangechain

import (
	"github.com/halprin/rangechain/generator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromSlice(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	expectedOutput := []interface{}{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	chain := FromSlice(inputSlice)

	slice, err := chain.Slice()
	assert.Equal(expectedOutput, slice)
	assert.Nil(err)
}

func TestFromArray(t *testing.T) {
	assert := assert.New(t)

	input := [...]string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	expectedOutput := []interface{}{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	chain := FromArray(input)

	slice, err := chain.Slice()
	assert.Equal(expectedOutput, slice)
	assert.Nil(err)
}

func TestFromChannel(t *testing.T) {
	assert := assert.New(t)

	innerInput := []string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	expectedOutput := []interface{}{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	input := createTestChannel(innerInput)
	chain := FromChannel(input)

	slice, err := chain.Slice()
	assert.Equal(expectedOutput, slice)
	assert.Nil(err)
}

func TestFromMap(t *testing.T) {
	assert := assert.New(t)

	key1 := "DogCow"
	value1 := 3
	key2 := "goes"
	value2 := 92
	key3 := "Moof!"
	value3 := 26

	input := map[string]int{
		key1: value1,
		key2: value2,
		key3: value3,
	}
	chain := FromMap(input)

	expectedOutput := []interface{}{
		generator.MapTuple{
			Key: key1,
			Value: value1,
		},
		generator.MapTuple{
			Key: key2,
			Value: value2,
		},
		generator.MapTuple{
			Key: key3,
			Value: value3,
		},
	}

	slice, err := chain.Slice()
	//not testing the order because we are not guaranteed the order in which a map is iterated over
	assert.ElementsMatch(expectedOutput, slice)
	assert.Nil(err)
}

func createTestChannel(stringSlice []string) chan string {
	stringChannel := make(chan string)

	go func() {
		for _, currentString := range stringSlice {
			stringChannel <- currentString
		}
		close(stringChannel)
	}()

	return stringChannel
}
