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
	slice, err := chain.Map(func(value interface{}) (interface{}, error) {
		intValue := value.(int)
		if intValue == errorValue {
			return intValue, errors.New("an example error yo")
		}

		return intValue, nil
	}).Skip(2).Slice()

	assert.Equal([]interface{}{8, 26}, slice)
	assert.Nil(err)
}