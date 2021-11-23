package rangechain

import (
	"github.com/halprin/rangechain/keyvalue"
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
	input := createTestStringChannel(innerInput)
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
		&testKeyValue{
			TheKey:   key1,
			TheValue: value1,
		},
		&testKeyValue{
			TheKey:   key2,
			TheValue: value2,
		},
		&testKeyValue{
			TheKey:   key3,
			TheValue: value3,
		},
	}

	slice, err := chain.Slice()
	//not testing the order because we are not guaranteed the order in which a map is iterated over
	assertEqualsBasedOnKeyValuerInterface(t, expectedOutput, slice)
	assert.Nil(err)
}

type testKeyValue struct {
	TheKey   interface{}
	TheValue interface{}
}

func (t *testKeyValue) Key() interface{} {
	return t.TheKey
}

func (t *testKeyValue) Value() interface{} {
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

func assertEqualsBasedOnKeyValuerInterface(t *testing.T, expected []interface{}, actual []interface{}) {
	assert := assert.New(t)

	assert.Len(actual, len(expected))

	for _, expectedValue := range expected {
		expectedKeyValuer, isType := expectedValue.(keyvalue.KeyValuer)
		if !isType {
			continue
		}
		keyToFind := expectedKeyValuer.Key()
		foundMatch := false

		for _, actualValue := range actual {
			actualKeyValuer, isType := actualValue.(keyvalue.KeyValuer)
			if !isType {
				continue
			}
			actualKey := actualKeyValuer.Key()

			if actualKey != keyToFind {
				continue
			}

			if expectedKeyValuer.Value() != actualKeyValuer.Value() {
				continue
			}

			foundMatch = true
		}

		assert.True(foundMatch)
	}
}
