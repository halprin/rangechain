package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterfaceSliceConversion(t *testing.T) {
	inputSlice := []string{"DogCow", "Moof", "Clarus"}
	actualSlice := InterfaceSlice(inputSlice)

	assert.ElementsMatch(t, inputSlice, actualSlice)
}

func TestInterfaceSlicePanicsForNonSlice(t *testing.T) {
	assert.Panics(t, func() {
		InterfaceSlice(3)
	})
}

func TestInterfaceSliceReturnsNil(t *testing.T) {
	var inputSlice []int

	actualSlice := InterfaceSlice(inputSlice)

	assert.Nil(t, actualSlice)
}
