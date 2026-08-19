package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	assert := assert.New(t)

	aSet := NewSet[string]()
	valueToCheck := "Moof!"

	assert.False(aSet.Contains(valueToCheck))

	aSet.Add(valueToCheck)

	assert.True(aSet.Contains(valueToCheck))
}
