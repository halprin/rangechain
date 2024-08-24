package generator

import (
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func TestFromSliceWithLastTimeError(t *testing.T) {
	assert := assert.New(t)

	generator := FromSlice([]interface{}{9})

	_, err := generator()

	assert.NoError(err)

	_, err = generator()

	assert.ErrorIs(err, Exhausted)
}

func TestFromSliceWithCorrectOrder(t *testing.T) {
	assert := assert.New(t)

	expectedFirstItem := 1
	expectedSecondItem := 26
	expectedThirdItem := 9
	theSlice := []interface{}{expectedFirstItem, expectedSecondItem, expectedThirdItem}
	generator := FromSlice(theSlice)

	actualFirstItem, err := generator()
	assert.Equal(expectedFirstItem, actualFirstItem)
	assert.NoError(err)

	actualSecondItem, err := generator()
	assert.Equal(expectedSecondItem, actualSecondItem)
	assert.NoError(err)

	actualThirdItem, err := generator()
	assert.Equal(actualThirdItem, actualThirdItem)
	assert.NoError(err)
}

func TestFromSliceEmpty(t *testing.T) {
	generator := FromSlice([]interface{}{})

	_, err := generator()

	assert.ErrorIs(t, err, Exhausted)
}

//func TestArrayDoesntWorkForSlice(t *testing.T) {
//	assert.Panics(t, func() {
//		FromSlice([...]interface{}{})
//	})
//}

func TestArrayNotPanicsGivenArray(t *testing.T) {
	assert.NotPanics(t, func() {
		FromArray([...]interface{}{})
	})
}

func TestArrayDoesntWorkGivenSlice(t *testing.T) {
	assert.Panics(t, func() {
		FromArray([]interface{}{})
	})
}

func TestChannelNotPanicsGivenChannel(t *testing.T) {
	assert.NotPanics(t, func() {
		FromChannel(createTestChannel(2))
	})
}

func TestChannelDoesntWorkGivenSlice(t *testing.T) {
	assert.Panics(t, func() {
		FromChannel([]interface{}{})
	})
}

func TestChannelEndsWithError(t *testing.T) {
	assert := assert.New(t)

	generator := FromChannel(createTestChannel(1))

	_, err := generator()
	assert.NoError(err)

	_, err = generator()
	assert.ErrorIs(err, Exhausted)
}

func TestChannelEndsWithErrorAndDoesntPanicAfterward(t *testing.T) {
	assert := assert.New(t)

	generator := FromChannel(createTestChannel(1))

	_, err := generator()
	assert.NoError(err)

	_, err = generator()
	assert.ErrorIs(err, Exhausted)

	assert.NotPanics(func() {
		_, err = generator()
		assert.ErrorIs(err, Exhausted)
	})
}

func TestChannelImmediatelyEnds(t *testing.T) {
	generator := FromChannel(createTestChannel(0))

	_, err := generator()

	assert.ErrorIs(t, err, Exhausted)
}

func TestMapNotPanicsGivenMap(t *testing.T) {
	assert.NotPanics(t, func() {
		FromMap(map[string]int{
			"DogCow": 1,
			"Moof":   1976,
		})
	})
}

func TestMapDoesntWorkGivenSlice(t *testing.T) {
	assert.Panics(t, func() {
		FromMap([]interface{}{})
	})
}

func TestMapEndsWithError(t *testing.T) {
	assert := assert.New(t)

	generator := FromMap(map[string]int{
		"DogCow": 3,
	})

	_, err := generator()
	assert.NoError(err)

	_, err = generator()
	assert.ErrorIs(err, Exhausted)
}

func TestFromIteratorWithLastTimeError(t *testing.T) {
	assert := assert.New(t)

	generator := FromIterator(slices.Values([]interface{}{9}))

	_, err := generator()

	assert.NoError(err)

	_, err = generator()

	assert.ErrorIs(err, Exhausted)
}

func TestFromIteratorEmpty(t *testing.T) {
	generator := FromIterator(slices.Values([]interface{}{}))

	_, err := generator()

	assert.ErrorIs(t, err, Exhausted)
}

func createTestChannel(size int) chan interface{} {
	intChannel := make(chan interface{})

	go func() {
		for i := 0; i < size; i++ {
			intChannel <- i
		}
		close(intChannel)
	}()

	return intChannel
}
