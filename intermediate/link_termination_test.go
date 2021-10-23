package intermediate


import (
	"errors"
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
	assert := assert.New(t)

	expectedSlice := []int{987, 8, 26}
	generation := generator.FromSlice(expectedSlice)
	link := NewLink(generation)

	var seenItems []interface{}
	valueChannel, errorChannel := link.Channel()
	for currentValue := range valueChannel {
		seenItems = append(seenItems, currentValue)
	}
	seenError := <-errorChannel

	assert.Equal(helper.InterfaceSlice(expectedSlice), seenItems)
	assert.Nil(seenError)
}

func TestChannelHasError(t *testing.T) {
	assert := assert.New(t)

	errorValue := 8
	inputSlice := []int{987, errorValue, 26}
	expectedError := errors.New("an example error yo")
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := NewLink(generation)

	var seenItems []interface{}
	valueChannel, errorChannel := link.Channel()
	//still range through the value channel to ensure we close the channel when an error is encountered
	for currentValue := range valueChannel {
		seenItems = append(seenItems, currentValue)
	}
	seenError := <-errorChannel

	assert.Equal(expectedError, seenError)
}

func TestForEach(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	var seenItems []interface{}
	forEachFunction := func(value interface{}) {
		seenItems = append(seenItems, value)
	}
	err := link.ForEach(forEachFunction)

	assert.ElementsMatch(inputSlice, seenItems)
	assert.Nil(err)
}

func TestForEachHasError(t *testing.T) {
	errorValue := 8
	inputSlice := []int{987, errorValue, 26}
	expectedError := errors.New("an example error yo")
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := NewLink(generation)

	forEachFunction := func(value interface{}) {}
	err := link.ForEach(forEachFunction)

	assert.Equal(t, expectedError, err)
}

func TestForEachParallel(t *testing.T) {
	assert := assert.New(t)

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
	err := link.ForEachParallel(forEachFunction)

	for {
		seenItemsLock.RLock()
		if len(seenItems) == len(expectedOutput) {
			seenItemsLock.RUnlock()
			break
		}
		seenItemsLock.RUnlock()
	}
	assert.Equal(expectedOutput, seenItems)
	assert.Nil(err)
}

func TestForEachParallelHasError(t *testing.T) {
	assert := assert.New(t)

	errorValue := 8
	inputSlice := []int{987, errorValue, 26}
	expectedError := errors.New("an example error yo")
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := NewLink(generation)

	forEachFunction := func(value interface{}) {}
	err := link.ForEachParallel(forEachFunction)

	assert.Equal(expectedError, err)
}

func TestCount(t *testing.T) {
	assert := assert.New(t)

	errorValue := 8
	inputSlice := []int{987, errorValue, 26}
	expectedError := errors.New("an example error yo")
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := NewLink(generation)

	actualCount, err := link.Count()

	assert.Equal(len(inputSlice), actualCount)
	assert.Equal(expectedError, err)
}

func TestCountWithErrorStillCounts(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	actualCount, err := link.Count()

	assert.Equal(len(inputSlice), actualCount)
	assert.Nil(err)
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
