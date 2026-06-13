package rangechain

import (
	"errors"
	"github.com/halprin/rangechain/internal/generator"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestSlice(t *testing.T) {
	assert := assert.New(t)
	expectedSlice := []int{987, 8, 26}
	generation := generator.FromSlice(expectedSlice)
	link := newLink(generation)

	actualSlice, err := link.Slice()

	assert.Equal(expectedSlice, actualSlice)
	assert.Nil(err)
}

func TestSliceWithErrorReturnsPartOfSlice(t *testing.T) {
	assert := assert.New(t)

	errorValue := 26
	inputSlice := []int{987, 8, errorValue}
	expectedError := errors.New("an example error yo")
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := newLink(generation)

	actualSlice, err := link.Slice()

	assert.Equal([]int{987, 8}, actualSlice)
	assert.Equal(expectedError, err)
}

func TestChannel(t *testing.T) {
	assert := assert.New(t)

	expectedSlice := []int{987, 8, 26}
	generation := generator.FromSlice(expectedSlice)
	link := newLink(generation)

	var seenItems []int
	valueChannel, errorChannel := link.Channel()
	for currentValue := range valueChannel {
		seenItems = append(seenItems, currentValue)
	}
	seenError := <-errorChannel

	assert.Equal(expectedSlice, seenItems)
	assert.Nil(seenError)
}

func TestChannelHasError(t *testing.T) {
	assert := assert.New(t)

	errorValue := 8
	inputSlice := []int{987, errorValue, 26}
	expectedError := errors.New("an example error yo")
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := newLink(generation)

	var seenItems []int
	valueChannel, errorChannel := link.Channel()
	for currentValue := range valueChannel {
		seenItems = append(seenItems, currentValue)
	}
	seenError := <-errorChannel

	assert.Equal(expectedError, seenError)
}

func TestIterator(t *testing.T) {
	assert := assert.New(t)

	expectedSlice := []int{987, 8, 26}
	generation := generator.FromSlice(expectedSlice)
	link := newLink(generation)

	var seenItems []int
	var seenError error
	for value, err := range link.Iterator() {
		if err != nil {
			seenError = err
			continue
		}
		seenItems = append(seenItems, value)
	}

	assert.Equal(expectedSlice, seenItems)
	assert.Nil(seenError)
}

func TestIteratorAfterChainMethods(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{1, 2, 3, 4, 5}
	link := FromSlice(inputSlice).
		Filter(func(value int) (bool, error) {
			return value%2 == 1, nil
		}).
		Map(func(value int) (int, error) {
			return value * 10, nil
		})

	var seenItems []int
	for value, err := range link.Iterator() {
		assert.Nil(err)
		seenItems = append(seenItems, value)
	}

	assert.Equal([]int{10, 30, 50}, seenItems)
}

func TestIteratorEarlyBreak(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987, 8, 26, 42, 100}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	var seenItems []int
	for value, err := range link.Iterator() {
		assert.Nil(err)
		seenItems = append(seenItems, value)
		if len(seenItems) == 2 {
			break
		}
	}

	assert.Equal([]int{987, 8}, seenItems)
}

func TestIteratorHasError(t *testing.T) {
	assert := assert.New(t)

	errorValue := 8
	inputSlice := []int{987, errorValue, 26}
	expectedError := errors.New("an example error yo")
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := newLink(generation)

	var seenItems []int
	var seenError error
	for value, err := range link.Iterator() {
		if err != nil {
			seenError = err
			continue
		}
		seenItems = append(seenItems, value)
	}

	assert.Equal([]int{987}, seenItems)
	assert.Equal(expectedError, seenError)
}

func TestForEach(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	var seenItems []int
	forEachFunction := func(value int) {
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
	link := newLink(generation)

	forEachFunction := func(value int) {}
	err := link.ForEach(forEachFunction)

	assert.Equal(t, expectedError, err)
}

func TestForEachParallel(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987, 8, 26}
	expectedOutput := map[int]bool{
		987: true,
		26:  true,
		8:   true,
	}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	seenItems := map[int]bool{}
	seenItemsLock := sync.RWMutex{}

	forEachFunction := func(value int) {
		seenItemsLock.Lock()
		seenItems[value] = true
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
	link := newLink(generation)

	forEachFunction := func(value int) {}
	err := link.ForEachParallel(forEachFunction)

	assert.Equal(expectedError, err)
}

func TestCount(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	actualCount, err := link.Count()

	assert.Equal(len(inputSlice), actualCount)
	assert.Nil(err)
}

func TestCountWithErrorStillCounts(t *testing.T) {
	assert := assert.New(t)

	errorValue := 8
	inputSlice := []int{987, errorValue, errorValue}
	expectedError := errors.New("an example error yo")
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := newLink(generation)

	actualCount, err := link.Count()

	assert.Equal(len(inputSlice), actualCount)
	assert.Equal(expectedError, err)
}

func TestCountWithFirstErrorReturns(t *testing.T) {
	assert := assert.New(t)

	errorValue := 8
	secondErrorValue := 26
	inputSlice := []int{987, errorValue, secondErrorValue}
	expectedError := errors.New("an example error yo")
	secondError := errors.New("another error")
	generation := wrapGeneratorWithError(createGeneratorWithError(inputSlice, errorValue, expectedError), secondErrorValue, secondError)
	link := newLink(generation)

	_, err := link.Count()

	assert.Equal(expectedError, err)
}

func TestFirst(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	actualFirst, err := link.First()

	assert.NotNil(actualFirst)
	assert.Equal(inputSlice[0], *actualFirst)
	assert.Nil(err)
}

func TestFirstWithEmptySlice(t *testing.T) {
	assert := assert.New(t)

	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

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
	link := newLink(generation)

	actualFirst, err := link.First()

	assert.Nil(actualFirst)
	assert.Equal(expectedError, err)
}

func TestLast(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	actualLast, err := link.Last()

	assert.NotNil(actualLast)
	assert.Equal(inputSlice[len(inputSlice)-1], *actualLast)
	assert.Nil(err)
}

func TestLastWithEmptySlice(t *testing.T) {
	assert := assert.New(t)

	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

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
	link := newLink(generation)

	_, err := link.Last()

	assert.Equal(expectedError, err)
}

func TestAllMatch(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{984, 8, 26}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	allMatchFunction := func(value int) (bool, error) {
		return value%2 == 0, nil
	}
	match, err := link.AllMatch(allMatchFunction)

	assert.True(match)
	assert.Nil(err)
}

func TestNotAllMatch(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{984, 7, 26}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	allMatchFunction := func(value int) (bool, error) {
		return value%2 == 0, nil
	}
	match, err := link.AllMatch(allMatchFunction)

	assert.False(match)
	assert.Nil(err)
}

func TestAllMatchWithEmptySlice(t *testing.T) {
	assert := assert.New(t)

	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	allMatchFunction := func(value int) (bool, error) {
		return value%2 == 0, nil
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
	link := newLink(generation)

	allMatchFunction := func(value int) (bool, error) {
		return value%2 == 0, nil
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
	link := newLink(generation)

	allMatchFunction := func(value int) (bool, error) {
		if value == errorValue {
			return true, expectedError
		}
		return value%2 == 0, nil
	}
	match, err := link.AllMatch(allMatchFunction)

	assert.False(match)
	assert.Equal(expectedError, err)
}

func TestAnyMatch(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{985, 3, 26}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	anyMatchFunction := func(value int) (bool, error) {
		return value%2 == 0, nil
	}
	match, err := link.AnyMatch(anyMatchFunction)

	assert.True(match)
	assert.Nil(err)
}

func TestNotAnyMatch(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{985, 7, 29}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	anyMatchFunction := func(value int) (bool, error) {
		return value%2 == 0, nil
	}
	match, err := link.AnyMatch(anyMatchFunction)

	assert.False(match)
	assert.Nil(err)
}

func TestAnyMatchWithEmptySlice(t *testing.T) {
	assert := assert.New(t)

	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	anyMatchFunction := func(value int) (bool, error) {
		return value%2 == 0, nil
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
	link := newLink(generation)

	anyMatchFunction := func(value int) (bool, error) {
		return value%2 == 0, nil
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
	link := newLink(generation)

	anyMatchFunction := func(value int) (bool, error) {
		if value == errorValue {
			return true, expectedError
		}
		return value%2 == 0, nil
	}
	match, err := link.AnyMatch(anyMatchFunction)

	assert.False(match)
	assert.Equal(expectedError, err)
}

func TestNoneMatch(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{985, 3, 27}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	noneMatchFunction := func(value int) (bool, error) {
		return value%2 == 0, nil
	}
	noneMatch, err := link.NoneMatch(noneMatchFunction)

	assert.True(noneMatch)
	assert.Nil(err)
}

func TestNotNoneMatch(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{985, 7, 28}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	noneMatchFunction := func(value int) (bool, error) {
		return value%2 == 0, nil
	}
	noneMatch, err := link.NoneMatch(noneMatchFunction)

	assert.False(noneMatch)
	assert.Nil(err)
}

func TestNoneMatchWithEmptySlice(t *testing.T) {
	assert := assert.New(t)

	var inputSlice []int
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	noneMatchFunction := func(value int) (bool, error) {
		return value%2 == 0, nil
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
	link := newLink(generation)

	noneMatchFunction := func(value int) (bool, error) {
		return value%2 == 0, nil
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
	link := newLink(generation)

	noneMatchFunction := func(value int) (bool, error) {
		if value == errorValue {
			return false, expectedError
		}
		return value%2 == 0, nil
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
	link := newLink(generation)

	reduceFunction := func(firstItem int, secondItem int) (int, error) {
		return firstItem * secondItem, nil
	}
	actualReducedValue, err := link.Reduce(reduceFunction)

	assert.NotNil(actualReducedValue)
	assert.Equal(expectedValue, *actualReducedValue)
	assert.Nil(err)
}

func TestReduceWithOneItem(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	reduceFunction := func(firstItem int, secondItem int) (int, error) {
		return firstItem * secondItem, nil
	}
	actualReducedValue, err := link.Reduce(reduceFunction)

	assert.NotNil(actualReducedValue)
	assert.Equal(inputSlice[0], *actualReducedValue)
	assert.Nil(err)
}

func TestReduceWithZeroItems(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	reduceFunction := func(firstItem int, secondItem int) (int, error) {
		return firstItem * secondItem, nil
	}
	actualReducedValue, err := link.Reduce(reduceFunction)

	assert.Nil(actualReducedValue)
	assert.Nil(err)
}

func TestReduceWithEarlierError(t *testing.T) {
	assert := assert.New(t)

	errorValue := 26
	inputSlice := []int{987, 8, errorValue}
	expectedError := errors.New("an example error yo")
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := newLink(generation)

	reduceFunction := func(firstItem int, secondItem int) (int, error) {
		return firstItem * secondItem, nil
	}
	_, err := link.Reduce(reduceFunction)

	assert.Equal(expectedError, err)
}

func TestReduceWithErrorInReduceFunction(t *testing.T) {
	assert := assert.New(t)

	errorValue := 26
	inputSlice := []int{987, 8, errorValue}
	expectedError := errors.New("an example error yo")
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	reduceFunction := func(firstItem int, secondItem int) (int, error) {
		if firstItem == errorValue || secondItem == errorValue {
			return 0, expectedError
		}
		return firstItem * secondItem, nil
	}
	_, err := link.Reduce(reduceFunction)

	assert.Equal(expectedError, err)
}

func TestReduceWithInitialValue(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987, 8, 26}
	inputInitialValue := 4
	expectedValue := inputInitialValue * inputSlice[0] * inputSlice[1] * inputSlice[2]
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	reduceFunction := func(firstItem int, secondItem int) (int, error) {
		return firstItem * secondItem, nil
	}
	actualReducedValue, err := link.ReduceWithInitialValue(reduceFunction, inputInitialValue)

	assert.Equal(expectedValue, actualReducedValue)
	assert.Nil(err)
}

func TestReduceWithInitialValueWithOneItem(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{987}
	inputInitialValue := 4
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	reduceFunction := func(firstItem int, secondItem int) (int, error) {
		return firstItem * secondItem, nil
	}
	actualReducedValue, err := link.ReduceWithInitialValue(reduceFunction, inputInitialValue)

	assert.Equal(inputInitialValue*inputSlice[0], actualReducedValue)
	assert.Nil(err)
}

func TestReduceWithInitialValueWithZeroItems(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{}
	inputInitialValue := 4
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	reduceFunction := func(firstItem int, secondItem int) (int, error) {
		return firstItem * secondItem, nil
	}
	actualReducedValue, err := link.ReduceWithInitialValue(reduceFunction, inputInitialValue)

	assert.Equal(inputInitialValue, actualReducedValue)
	assert.Nil(err)
}

func TestReduceWithInitialValueWithEarlierError(t *testing.T) {
	assert := assert.New(t)

	errorValue := 8
	inputSlice := []int{987, errorValue, 26}
	expectedError := errors.New("an example error yo")
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := newLink(generation)

	reduceFunction := func(firstItem int, secondItem int) (int, error) {
		return firstItem * secondItem, nil
	}
	_, err := link.ReduceWithInitialValue(reduceFunction, 4)

	assert.Equal(expectedError, err)
}

func TestReduceWithInitialValueWithErrorInReduceFunction(t *testing.T) {
	assert := assert.New(t)

	errorValue := 8
	inputSlice := []int{987, errorValue, 26}
	expectedError := errors.New("an example error yo")
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	reduceFunction := func(firstItem int, secondItem int) (int, error) {
		if firstItem == errorValue || secondItem == errorValue {
			return 0, expectedError
		}
		return firstItem * secondItem, nil
	}
	_, err := link.ReduceWithInitialValue(reduceFunction, 4)

	assert.Equal(expectedError, err)
}

func wrapGeneratorWithError(generation func() (int, error), valueToErrorOn int, errorToReturn error) func() (int, error) {
	return func() (int, error) {
		value, err := generation()

		if err != nil {
			return value, err
		} else if value == valueToErrorOn {
			return 0, errorToReturn
		}

		return value, err
	}
}
