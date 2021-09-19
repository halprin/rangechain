package intermediate


import (
	"github.com/halprin/slice-chain/generator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlice(t *testing.T) {
	expectedSlice := []int{987, 8, 26}
	generation := generator.FromSlice(expectedSlice)
	link := NewLink(generation)

	actualSlice := link.Slice()

	assert.Equal(t, expectedSlice, actualSlice)
}

func TestForEach(t *testing.T) {
	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	var seenItems []int
	forEachFunction := func(value int) {
		seenItems = append(seenItems, value)
	}
	link.ForEach(forEachFunction)

	assert.ElementsMatch(t, inputSlice, seenItems)
}

func TestCount(t *testing.T) {
	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	actualCount := link.Count()

	assert.Equal(t, len(inputSlice), actualCount)
}

func TestFirst(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	actualFirst := link.First()

	assert.NotNil(actualFirst)
	assert.Equal(inputSlice[0], *actualFirst)
}

func TestFirstWithEmptySlice(t *testing.T) {
	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	actualFirst := link.First()

	assert.Nil(t, actualFirst)
}

func TestLast(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	actualLast := link.Last()

	assert.NotNil(actualLast)
	assert.Equal(inputSlice[len(inputSlice) - 1], *actualLast)
}

func TestLastWithEmptySlice(t *testing.T) {
	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	actualLast := link.Last()

	assert.Nil(t, actualLast)
}

func TestAllMatch(t *testing.T) {
	inputSlice := []int{984, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	allMatchFunction := func(value int) bool {
		return value % 2 == 0  //even means true
	}
	match := link.AllMatch(allMatchFunction)

	assert.True(t, match)
}

func TestNotAllMatch(t *testing.T) {
	inputSlice := []int{984, 7, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	allMatchFunction := func(value int) bool {
		return value % 2 == 0  //even means true
	}
	match := link.AllMatch(allMatchFunction)

	assert.False(t, match)
}

func TestAllMatchWithEmptySlice(t *testing.T) {
	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	allMatchFunction := func(value int) bool {
		return value % 2 == 0  //even means true
	}
	match := link.AllMatch(allMatchFunction)

	assert.True(t, match)
}

func TestAnyMatch(t *testing.T) {
	inputSlice := []int{985, 3, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	anyMatchFunction := func(value int) bool {
		return value % 2 == 0  //even means true
	}
	match := link.AnyMatch(anyMatchFunction)

	assert.True(t, match)
}

func TestNotAnyMatch(t *testing.T) {
	inputSlice := []int{985, 7, 29}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	anyMatchFunction := func(value int) bool {
		return value % 2 == 0  //even means true
	}
	match := link.AnyMatch(anyMatchFunction)

	assert.False(t, match)
}

func TestAnyMatchWithEmptySlice(t *testing.T) {
	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	anyMatchFunction := func(value int) bool {
		return value % 2 == 0  //even means true
	}
	match := link.AnyMatch(anyMatchFunction)

	assert.False(t, match)
}

func TestNoneMatch(t *testing.T) {
	inputSlice := []int{985, 3, 27}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	noneMatchFunction := func(value int) bool {
		return value % 2 == 0  //even means true
	}
	match := link.NoneMatch(noneMatchFunction)

	assert.True(t, match)
}

func TestNotNoneMatch(t *testing.T) {
	inputSlice := []int{985, 7, 28}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	noneMatchFunction := func(value int) bool {
		return value % 2 == 0  //even means true
	}
	match := link.NoneMatch(noneMatchFunction)

	assert.False(t, match)
}

func TestNoneMatchWithEmptySlice(t *testing.T) {
	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	noneMatchFunction := func(value int) bool {
		return value % 2 == 0  //even means true
	}
	match := link.NoneMatch(noneMatchFunction)

	assert.True(t, match)
}