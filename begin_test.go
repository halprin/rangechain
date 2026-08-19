package rangechain

import (
	"maps"
	"slices"
	"testing"

	"github.com/halprin/rangechain/keyvalue"
	"github.com/stretchr/testify/assert"
)

func TestFromSlice(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	chain := FromSlice(inputSlice)

	slice, err := chain.Slice()
	assert.Equal(inputSlice, slice)
	assert.Nil(err)
}

func TestFromChannel(t *testing.T) {
	assert := assert.New(t)

	innerInput := []string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	input := createTestStringChannel(innerInput)
	chain := FromChannel((<-chan string)(input))

	slice, err := chain.Slice()
	assert.Equal(innerInput, slice)
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

	expectedOutput := []keyvalue.KeyValuer[string, int]{
		&testKeyValue[string, int]{TheKey: key1, TheValue: value1},
		&testKeyValue[string, int]{TheKey: key2, TheValue: value2},
		&testKeyValue[string, int]{TheKey: key3, TheValue: value3},
	}

	slice, err := chain.Slice()
	assertEqualsBasedOnKeyValuerInterface(t, expectedOutput, slice)
	assert.Nil(err)
}

func TestFromSeq2(t *testing.T) {
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
	chain := FromSeq2(maps.All(input))

	expectedOutput := []keyvalue.KeyValuer[string, int]{
		&testKeyValue[string, int]{TheKey: key1, TheValue: value1},
		&testKeyValue[string, int]{TheKey: key2, TheValue: value2},
		&testKeyValue[string, int]{TheKey: key3, TheValue: value3},
	}

	slice, err := chain.Slice()
	assertEqualsBasedOnKeyValuerInterface(t, expectedOutput, slice)
	assert.Nil(err)
}

func TestFromIterator(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	inputIterator := slices.Values(inputSlice)
	chain := FromIterator(inputIterator)

	slice, err := chain.Slice()
	assert.Equal(inputSlice, slice)
	assert.Nil(err)
}

type testKeyValue[K, V any] struct {
	TheKey   K
	TheValue V
}

func (t *testKeyValue[K, V]) Key() K {
	return t.TheKey
}

func (t *testKeyValue[K, V]) Value() V {
	return t.TheValue
}

func createTestStringChannel(stringSlice []string) chan string {
	stringChannel := make(chan string)

	go func() {
		for _, currentString := range stringSlice {
			stringChannel <- currentString
		}
		close(stringChannel)
	}()

	return stringChannel
}

func assertEqualsBasedOnKeyValuerInterface[K comparable, V comparable](t *testing.T, expected []keyvalue.KeyValuer[K, V], actual []keyvalue.KeyValuer[K, V]) {
	assert := assert.New(t)

	assert.Len(actual, len(expected))

	for _, expectedValue := range expected {
		keyToFind := expectedValue.Key()
		foundMatch := false

		for _, actualValue := range actual {
			actualKey := actualValue.Key()

			if actualKey != keyToFind {
				continue
			}

			if expectedValue.Value() != actualValue.Value() {
				continue
			}

			foundMatch = true
		}

		assert.True(foundMatch)
	}
}
