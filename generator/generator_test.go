package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromSliceWithLastTimePanic(t *testing.T) {
	generator := FromSlice([]int{9})

	assert.NotPanics(t, func() {
		generator()
	})

	assert.Panics(t, func() {
		generator()
	})
}

func TestFromSliceWithCorrectOrder(t *testing.T) {
	assert := assert.New(t)

	expectedFirstItem := 1
	expectedSecondItem := 26
	expectedThirdItem := 9
	theSlice := []int{expectedFirstItem, expectedSecondItem, expectedThirdItem}
	generator := FromSlice(theSlice)

	actualFirstItem := generator()
	assert.Equal(expectedFirstItem, actualFirstItem)

	actualSecondItem := generator()
	assert.Equal(expectedSecondItem, actualSecondItem)

	actualThirdItem := generator()
	assert.Equal(actualThirdItem, actualThirdItem)
}

func TestFromSliceEmpty(t *testing.T) {
	iterator := FromSlice([]int{})

	assert.Panics(t, func() {
		iterator()
	})
}
