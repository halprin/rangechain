package intermediate


import (
	"github.com/halprin/rangechain/generator"
	"github.com/halprin/rangechain/helper"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestSlice(t *testing.T) {
	assert := assert.New(t)
	expectedSlice := []int{987, 8, 26}
	generation := generator.FromSlice(expectedSlice)
	link := NewLink(generation)

	actualSlice, err := link.Slice()

	assert.Equal(helper.InterfaceSlice(expectedSlice), actualSlice)
	assert.Nil(err)
}

func TestChannel(t *testing.T) {
	expectedSlice := []int{987, 8, 26}
	generation := generator.FromSlice(expectedSlice)
	link := NewLink(generation)

	var seenItems []interface{}
	for currentValue := range link.Channel() {
		seenItems = append(seenItems, currentValue)
	}

	assert.Equal(t, helper.InterfaceSlice(expectedSlice), seenItems)
}

func TestForEach(t *testing.T) {
	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	var seenItems []interface{}
	forEachFunction := func(value interface{}) {
		seenItems = append(seenItems, value)
	}
	link.ForEach(forEachFunction)

	assert.ElementsMatch(t, inputSlice, seenItems)
}

func TestForEachParallel(t *testing.T) {
	inputSlice := []int{987, 8, 26}
	expectedOutput := map[int]bool{
		987: true,
		26: true,
		8: true,
	}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	seenItems := map[int]bool{}
	seenItemsLock := sync.RWMutex{}

	forEachFunction := func(value interface{}) {
		actualInt := value.(int)
		seenItemsLock.Lock()
		seenItems[actualInt] = true
		seenItemsLock.Unlock()
	}
	link.ForEachParallel(forEachFunction)

	for {
		seenItemsLock.RLock()
		if len(seenItems) == len(expectedOutput) {
			seenItemsLock.RUnlock()
			break
		}
		seenItemsLock.RUnlock()
	}
	assert.Equal(t, expectedOutput, seenItems)

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

	allMatchFunction := func(value interface{}) bool {
		intValue := value.(int)
		return intValue % 2 == 0  //even means true
	}
	match := link.AllMatch(allMatchFunction)

	assert.True(t, match)
}

func TestNotAllMatch(t *testing.T) {
	inputSlice := []int{984, 7, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	allMatchFunction := func(value interface{}) bool {
		intValue := value.(int)
		return intValue % 2 == 0  //even means true
	}
	match := link.AllMatch(allMatchFunction)

	assert.False(t, match)
}

func TestAllMatchWithEmptySlice(t *testing.T) {
	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	allMatchFunction := func(value interface{}) bool {
		intValue := value.(int)
		return intValue % 2 == 0  //even means true
	}
	match := link.AllMatch(allMatchFunction)

	assert.True(t, match)
}

func TestAnyMatch(t *testing.T) {
	inputSlice := []int{985, 3, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	anyMatchFunction := func(value interface{}) bool {
		intValue := value.(int)
		return intValue % 2 == 0  //even means true
	}
	match := link.AnyMatch(anyMatchFunction)

	assert.True(t, match)
}

func TestNotAnyMatch(t *testing.T) {
	inputSlice := []int{985, 7, 29}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	anyMatchFunction := func(value interface{}) bool {
		intValue := value.(int)
		return intValue % 2 == 0  //even means true
	}
	match := link.AnyMatch(anyMatchFunction)

	assert.False(t, match)
}

func TestAnyMatchWithEmptySlice(t *testing.T) {
	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	anyMatchFunction := func(value interface{}) bool {
		intValue := value.(int)
		return intValue % 2 == 0  //even means true
	}
	match := link.AnyMatch(anyMatchFunction)

	assert.False(t, match)
}

func TestNoneMatch(t *testing.T) {
	inputSlice := []int{985, 3, 27}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	noneMatchFunction := func(value interface{}) bool {
		intValue := value.(int)
		return intValue % 2 == 0  //even means true
	}
	match := link.NoneMatch(noneMatchFunction)

	assert.True(t, match)
}

func TestNotNoneMatch(t *testing.T) {
	inputSlice := []int{985, 7, 28}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	noneMatchFunction := func(value interface{}) bool {
		intValue := value.(int)
		return intValue % 2 == 0  //even means true
	}
	match := link.NoneMatch(noneMatchFunction)

	assert.False(t, match)
}

func TestNoneMatchWithEmptySlice(t *testing.T) {
	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	noneMatchFunction := func(value interface{}) bool {
		intValue := value.(int)
		return intValue % 2 == 0  //even means true
	}
	match := link.NoneMatch(noneMatchFunction)

	assert.True(t, match)
}

func TestReduce(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987, 8, 26}
	expectedValue := inputSlice[0] * inputSlice[1] * inputSlice[2]
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	reduceFunction := func(firstItem interface{}, secondItem interface{}) interface{} {
		firstIntItem := firstItem.(int)
		secondIntItem := secondItem.(int)

		return firstIntItem * secondIntItem
	}
	actualReducedValue := link.Reduce(reduceFunction)

	assert.NotNil(actualReducedValue)
	assert.Equal(expectedValue, *actualReducedValue)
}

func TestReduceWithOneItem(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	reduceFunction := func(firstItem interface{}, secondItem interface{}) interface{} {
		firstIntItem := firstItem.(int)
		secondIntItem := secondItem.(int)

		return firstIntItem * secondIntItem
	}
	actualReducedValue := link.Reduce(reduceFunction)

	assert.NotNil(actualReducedValue)
	assert.Equal(inputSlice[0], *actualReducedValue)
}

func TestReduceWithZeroItems(t *testing.T) {
	inputSlice := []int{}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	reduceFunction := func(firstItem interface{}, secondItem interface{}) interface{} {
		firstIntItem := firstItem.(int)
		secondIntItem := secondItem.(int)

		return firstIntItem * secondIntItem
	}
	actualReducedValue := link.Reduce(reduceFunction)

	assert.Nil(t, actualReducedValue)
}

func TestReduceWithInitialValue(t *testing.T) {
	inputSlice := []int{987, 8, 26}
	inputInitialValue := 4
	expectedValue := inputInitialValue * inputSlice[0] * inputSlice[1] * inputSlice[2]
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	reduceFunction := func(firstItem interface{}, secondItem interface{}) interface{} {
		firstIntItem := firstItem.(int)
		secondIntItem := secondItem.(int)

		return firstIntItem * secondIntItem
	}
	actualReducedValue := link.ReduceWithInitialValue(reduceFunction, inputInitialValue)

	assert.Equal(t, expectedValue, actualReducedValue)
}

func TestReduceWithInitialValueWithOneItem(t *testing.T) {
	inputSlice := []int{987}
	inputInitialValue := 4
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	reduceFunction := func(firstItem interface{}, secondItem interface{}) interface{} {
		firstIntItem := firstItem.(int)
		secondIntItem := secondItem.(int)

		return firstIntItem * secondIntItem
	}
	actualReducedValue := link.ReduceWithInitialValue(reduceFunction, inputInitialValue)

	assert.Equal(t, inputInitialValue * inputSlice[0], actualReducedValue)
}

func TestReduceWithInitialValueWithZeroItems(t *testing.T) {
	inputSlice := []int{}
	inputInitialValue := 4
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	reduceFunction := func(firstItem interface{}, secondItem interface{}) interface{} {
		firstIntItem := firstItem.(int)
		secondIntItem := secondItem.(int)

		return firstIntItem * secondIntItem
	}
	actualReducedValue := link.ReduceWithInitialValue(reduceFunction, inputInitialValue)

	assert.Equal(t, inputInitialValue, actualReducedValue)
}
