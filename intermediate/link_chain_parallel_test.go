package intermediate

import (
	"github.com/halprin/slice-chain/generator"
	"github.com/halprin/slice-chain/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapParallel(t *testing.T) {
	inputSlice := []string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	var expectedOutput []interface{}
	for _, stringValue := range inputSlice {
		expectedOutput = append(expectedOutput, len(stringValue))
	}

	generation := generator.FromSlice(helper.InterfaceSlice(inputSlice))
	link := NewLink(generation)

	mapFunction := func(value interface{}) interface{} {
		stringValue := value.(string)
		return len(stringValue)
	}

	actualSlice := link.MapParallel(mapFunction).Slice()

	assert.Equal(t, expectedOutput, actualSlice)
}
