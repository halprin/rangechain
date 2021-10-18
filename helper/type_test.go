package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var slice = []int{1, 2, 3}
var array = [...]int{1, 2, 3}
var channel = make(chan interface{})

func TestSlice(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsSlice(slice))
	assert.False(IsSlice(array))
	assert.False(IsSlice(channel))
}

func TestArray(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsArray(array))
	assert.False(IsArray(slice))
	assert.False(IsArray(channel))
}

func TestChannel(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsChannel(channel))
	assert.False(IsChannel(slice))
	assert.False(IsChannel(array))
}
