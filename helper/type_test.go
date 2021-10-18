package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var slice = []int{1, 2, 3}
var array = [3]int{1, 2, 3}

func TestSlice(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsSlice(slice))
	assert.False(IsSlice(array))
}

func TestArray(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsArray(array))
	assert.False(IsArray(slice))
}
