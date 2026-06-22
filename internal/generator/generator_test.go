package generator

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromSliceWithLastTimeError(t *testing.T) {
	assert := assert.New(t)

	gen := FromSlice([]int{9})

	_, err := gen()

	assert.NoError(err)

	_, err = gen()

	assert.ErrorIs(err, Exhausted)
}

func TestFromSliceWithCorrectOrder(t *testing.T) {
	assert := assert.New(t)

	expectedFirstItem := 1
	expectedSecondItem := 26
	expectedThirdItem := 9
	theSlice := []int{expectedFirstItem, expectedSecondItem, expectedThirdItem}
	gen := FromSlice(theSlice)

	actualFirstItem, err := gen()
	assert.Equal(expectedFirstItem, actualFirstItem)
	assert.NoError(err)

	actualSecondItem, err := gen()
	assert.Equal(expectedSecondItem, actualSecondItem)
	assert.NoError(err)

	actualThirdItem, err := gen()
	assert.Equal(expectedThirdItem, actualThirdItem)
	assert.NoError(err)
}

func TestFromSliceEmpty(t *testing.T) {
	gen := FromSlice([]int{})

	_, err := gen()

	assert.ErrorIs(t, err, Exhausted)
}

func TestChannelEndsWithError(t *testing.T) {
	assert := assert.New(t)

	gen := FromChannel(createTestChannel(1))

	_, err := gen()
	assert.NoError(err)

	_, err = gen()
	assert.ErrorIs(err, Exhausted)
}

func TestChannelEndsWithErrorAndDoesntPanicAfterward(t *testing.T) {
	assert := assert.New(t)

	gen := FromChannel(createTestChannel(1))

	_, err := gen()
	assert.NoError(err)

	_, err = gen()
	assert.ErrorIs(err, Exhausted)

	assert.NotPanics(func() {
		_, err = gen()
		assert.ErrorIs(err, Exhausted)
	})
}

func TestChannelImmediatelyEnds(t *testing.T) {
	gen := FromChannel(createTestChannel(0))

	_, err := gen()

	assert.ErrorIs(t, err, Exhausted)
}

func TestMapEndsWithError(t *testing.T) {
	assert := assert.New(t)

	gen := FromMap(map[string]int{
		"DogCow": 3,
	})

	_, err := gen()
	assert.NoError(err)

	_, err = gen()
	assert.ErrorIs(err, Exhausted)
}

func TestFromIteratorWithLastTimeError(t *testing.T) {
	assert := assert.New(t)

	gen := FromIterator(slices.Values([]int{9}))

	_, err := gen()

	assert.NoError(err)

	_, err = gen()

	assert.ErrorIs(err, Exhausted)
}

func TestFromIteratorEmpty(t *testing.T) {
	gen := FromIterator(slices.Values([]int{}))

	_, err := gen()

	assert.ErrorIs(t, err, Exhausted)
}

func createTestChannel(size int) chan int {
	intChannel := make(chan int)

	go func() {
		for i := 0; i < size; i++ {
			intChannel <- i
		}
		close(intChannel)
	}()

	return intChannel
}
