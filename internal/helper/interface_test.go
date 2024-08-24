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

func TestInterfaceMapConversion(t *testing.T) {
	inputMap := map[int]string{3: "DogCow", 26: "Moof", 1337: "Clarus"}
	actualMap := InterfaceMap(inputMap)

	assert.Len(t, actualMap, len(inputMap))

	for key, value := range inputMap {
		assert.Equal(t, value, actualMap[key])
	}
}

func TestInterfaceMapPanicsForNonSlice(t *testing.T) {
	assert.Panics(t, func() {
		InterfaceMap(3)
	})
}

func TestInterfaceMapReturnsNil(t *testing.T) {
	var inputMap map[int]string

	actualMap := InterfaceMap(inputMap)

	assert.Nil(t, actualMap)
}
