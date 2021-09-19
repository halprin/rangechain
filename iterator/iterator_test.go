package iterator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromSliceWithLastTimePanic(t *testing.T) {
	iterator := FromSlice([]int{9})

	assert.NotPanics(t, func() {
		iterator()
	})

	assert.Panics(t, func() {
		iterator()
	})
}

func TestFromSliceWithCorrectOrder(t *testing.T) {
	assert := assert.New(t)
	expectedFirstItem := 1
	expectedSecondItem := 26
	expectedThirdItem := 9
	theSlice := []int{expectedFirstItem, expectedSecondItem, expectedThirdItem}
	iterator := FromSlice(theSlice)

	actualFirstItem := iterator()
	assert.Equal(expectedFirstItem, actualFirstItem)

	actualSecondItem := iterator()
	assert.Equal(expectedSecondItem, actualSecondItem)

	actualThirdItem := iterator()
	assert.Equal(actualThirdItem, actualThirdItem)
}

func TestFromSliceEmpty(t *testing.T) {
	iterator := FromSlice([]int{})

	assert.Panics(t, func() {
		iterator()
	})
}
