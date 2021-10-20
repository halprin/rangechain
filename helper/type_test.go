package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var slice = []int{1, 2, 3}
var array = [...]int{1, 2, 3}
var channel = make(chan interface{})
var aMap = map[string]int{
	"a": 6,
	"DogCow": 1,
	"Moof": 1337,
}

func TestSlice(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsSlice(slice))
	assert.False(IsSlice(array))
	assert.False(IsSlice(channel))
	assert.False(IsSlice(aMap))
}

func TestArray(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsArray(array))
	assert.False(IsArray(slice))
	assert.False(IsArray(channel))
	assert.False(IsArray(aMap))
}

func TestChannel(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsChannel(channel))
	assert.False(IsChannel(slice))
	assert.False(IsChannel(array))
	assert.False(IsChannel(aMap))
}

func TestMap(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsMap(aMap))
	assert.False(IsMap(slice))
	assert.False(IsMap(array))
	assert.False(IsMap(channel))
}
