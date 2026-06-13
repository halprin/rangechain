package rangechain

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSkipSkipsErrors(t *testing.T) {
	assert := assert.New(t)

	errorValue := 987
	inputSlice := []int{errorValue, errorValue, 8, 26}
	chain := FromSlice(inputSlice)
	slice, err := chain.Map(func(value int) (int, error) {
		if value == errorValue {
			return value, errors.New("an example error yo")
		}

		return value, nil
	}).Skip(2).Slice()

	assert.Equal([]int{8, 26}, slice)
	assert.Nil(err)
}
