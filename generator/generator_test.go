package generator

import (
	"github.com/stretchr/testify/assert"
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

func TestArrayDoesntWorkForSlice(t *testing.T) {
	assert.Panics(t, func() {
		FromSlice([1]interface{}{})
	})
}

func TestArray(t *testing.T) {
	assert.NotPanics(t, func() {
		FromArray([0]interface{}{})
	})
}

func TestSliceDoesntWorkForArray(t *testing.T) {
	assert.Panics(t, func() {
		FromArray([]interface{}{})
	})
}
