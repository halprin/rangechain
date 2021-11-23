package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet(t *testing.T) {
	assert := assert.New(t)

	aSet := NewSet()
	valueToCheck := "Moof!"

	assert.False(aSet.Contains(valueToCheck))

	aSet.Add(valueToCheck)

	assert.True(aSet.Contains(valueToCheck))
}
