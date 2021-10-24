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

	actualFirst, err := link.First()

	assert.NotNil(actualFirst)
	assert.Equal(inputSlice[0], *actualFirst)
	assert.Nil(err)
}

func TestFirstWithEmptySlice(t *testing.T) {
	assert := assert.New(t)

	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	actualFirst, err := link.First()

	assert.Nil(actualFirst)
	assert.Nil(err)
}

func TestFirstHasError(t *testing.T) {
	assert := assert.New(t)

	errorValue := 987
	inputSlice := []int{errorValue, 8, 26}
	expectedError := errors.New("an example error yo")
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := NewLink(generation)

	actualFirst, err := link.First()

	assert.Nil(actualFirst)
	assert.Equal(expectedError, err)
}

func TestLast(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	actualLast, err := link.Last()

	assert.NotNil(actualLast)
	assert.Equal(inputSlice[len(inputSlice) - 1], *actualLast)
	assert.Nil(err)
}

func TestLastWithEmptySlice(t *testing.T) {
	assert := assert.New(t)

	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	actualLast, err := link.Last()

	assert.Nil(actualLast)
	assert.Nil(err)
}

func TestLastHasError(t *testing.T) {
	assert := assert.New(t)

	errorValue := 26
	inputSlice := []int{987, 8, errorValue}
	expectedError := errors.New("an example error yo")
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := NewLink(generation)

	_, err := link.Last()

	assert.Equal(expectedError, err)
}

func TestAllMatch(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{984, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	allMatchFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		return intValue % 2 == 0, nil  //even means true
	}
	match, err := link.AllMatch(allMatchFunction)

	assert.True(match)
	assert.Nil(err)
}

func TestNotAllMatch(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{984, 7, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	allMatchFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		return intValue % 2 == 0, nil  //even means true
	}
	match, err := link.AllMatch(allMatchFunction)

	assert.False(match)
	assert.Nil(err)
}

func TestAllMatchWithEmptySlice(t *testing.T) {
	assert := assert.New(t)

	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	allMatchFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		return intValue % 2 == 0, nil  //even means true
	}
	match, err := link.AllMatch(allMatchFunction)

	assert.True(match)
	assert.Nil(err)
}

func TestAllMatchWithEarlierError(t *testing.T) {
	assert := assert.New(t)

	errorValue := 8
	inputSlice := []int{984, errorValue, 26}
	expectedError := errors.New("an example error yo")
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := NewLink(generation)

	allMatchFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		return intValue % 2 == 0, nil  //even means true
	}
	match, err := link.AllMatch(allMatchFunction)

	assert.False(match)
	assert.Equal(expectedError, err)
}

func TestAllMatchWithErrorInMatchFunction(t *testing.T) {
	assert := assert.New(t)

	errorValue := 8
	inputSlice := []int{984, errorValue, 26}
	expectedError := errors.New("an example error yo")
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	allMatchFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		if intValue == errorValue {
			return true, expectedError
		}
		return intValue % 2 == 0, nil  //even means true
	}
	match, err := link.AllMatch(allMatchFunction)

	assert.False(match)
	assert.Equal(expectedError, err)
}

func TestAnyMatch(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{985, 3, 26}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	anyMatchFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		return intValue % 2 == 0, nil  //even means true
	}
	match, err := link.AnyMatch(anyMatchFunction)

	assert.True(match)
	assert.Nil(err)
}

func TestNotAnyMatch(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{985, 7, 29}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	anyMatchFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		return intValue % 2 == 0, nil  //even means true
	}
	match, err := link.AnyMatch(anyMatchFunction)

	assert.False(match)
	assert.Nil(err)
}

func TestAnyMatchWithEmptySlice(t *testing.T) {
	assert := assert.New(t)

	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	anyMatchFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		return intValue % 2 == 0, nil  //even means true
	}
	match, err := link.AnyMatch(anyMatchFunction)

	assert.False(match)
	assert.Nil(err)
}

func TestAnyMatchWithEarlierError(t *testing.T) {
	assert := assert.New(t)

	errorValue := 26
	inputSlice := []int{985, 3, errorValue}
	expectedError := errors.New("an example error yo")
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := NewLink(generation)

	anyMatchFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		return intValue % 2 == 0, nil  //even means true
	}
	match, err := link.AnyMatch(anyMatchFunction)

	assert.False(match)
	assert.Equal(expectedError, err)
}

func TestAnyMatchWithErrorInMatchFunction(t *testing.T) {
	assert := assert.New(t)

	errorValue := 26
	inputSlice := []int{985, 3, errorValue}
	expectedError := errors.New("an example error yo")
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	anyMatchFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		if intValue == errorValue {
			return true, expectedError
		}
		return intValue % 2 == 0, nil  //even means true
	}
	match, err := link.AnyMatch(anyMatchFunction)

	assert.False(match)
	assert.Equal(expectedError, err)
}

func TestNoneMatch(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{985, 3, 27}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	noneMatchFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		return intValue % 2 == 0, nil  //even means true
	}
	noneMatch, err := link.NoneMatch(noneMatchFunction)

	assert.True(noneMatch)
	assert.Nil(err)
}

func TestNotNoneMatch(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{985, 7, 28}
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	noneMatchFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		return intValue % 2 == 0, nil  //even means true
	}
	noneMatch, err := link.NoneMatch(noneMatchFunction)

	assert.False(noneMatch)
	assert.Nil(err)
}

func TestNoneMatchWithEmptySlice(t *testing.T) {
	assert := assert.New(t)

	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	noneMatchFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		return intValue % 2 == 0, nil  //even means true
	}
	noneMatch, err := link.NoneMatch(noneMatchFunction)

	assert.True(noneMatch)
	assert.Nil(err)
}

func TestNoneMatchWithEarlierError(t *testing.T) {
	assert := assert.New(t)

	errorValue := 3
	inputSlice := []int{985, errorValue, 27}
	expectedError := errors.New("an example error yo")
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := NewLink(generation)

	noneMatchFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		return intValue % 2 == 0, nil  //even means true
	}
	noneMatch, err := link.NoneMatch(noneMatchFunction)

	assert.True(noneMatch)
	assert.Equal(expectedError, err)
}

func TestNoneMatchWithErrorInMatchFunction(t *testing.T) {
	assert := assert.New(t)

	errorValue := 3
	inputSlice := []int{985, errorValue, 27}
	expectedError := errors.New("an example error yo")
	generation := generator.FromSlice(inputSlice)
	link := NewLink(generation)

	noneMatchFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		if intValue == errorValue {
			return false, expectedError
		}
		return intValue % 2 == 0, nil  //even means true
	}
	noneMatch, err := link.NoneMatch(noneMatchFunction)

	assert.True(noneMatch)
	assert.Equal(expectedError, err)
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
