package rangechain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var isCheckSlice = []int{1, 2, 3}
var isCheckArray = [...]int{1, 2, 3}
var isCheckChannel = make(chan any)
var isCheckMap = map[string]int{
	"a":      6,
	"DogCow": 1,
	"Moof":   1337,
}

func TestIsSlice(t *testing.T) {
	assert := assert.New(t)

	assert.True(isSlice(isCheckSlice))
	assert.False(isSlice(isCheckArray))
	assert.False(isSlice(isCheckChannel))
	assert.False(isSlice(isCheckMap))
}

func TestIsArray(t *testing.T) {
	assert := assert.New(t)

	assert.True(isArray(isCheckArray))
	assert.False(isArray(isCheckSlice))
	assert.False(isArray(isCheckChannel))
	assert.False(isArray(isCheckMap))
}

func TestIsChannel(t *testing.T) {
	assert := assert.New(t)

	assert.True(isChannel(isCheckChannel))
	assert.False(isChannel(isCheckSlice))
	assert.False(isChannel(isCheckArray))
	assert.False(isChannel(isCheckMap))
}

func TestIsMap(t *testing.T) {
	assert := assert.New(t)

	assert.True(isMap(isCheckMap))
	assert.False(isMap(isCheckSlice))
	assert.False(isMap(isCheckArray))
	assert.False(isMap(isCheckChannel))
}
