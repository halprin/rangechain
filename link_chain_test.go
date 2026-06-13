package rangechain

import (
	"errors"
	"github.com/halprin/rangechain/internal/generator"
	"github.com/halprin/rangechain/keyvalue"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	var expectedOutput []int
	for _, stringValue := range inputSlice {
		expectedOutput = append(expectedOutput, len(stringValue))
	}

	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	mapFunction := func(value string) (int, error) {
		return len(value), nil
	}

	actualSlice, err := link.Map(mapFunction).Slice()

	assert.Equal(expectedOutput, actualSlice)
	assert.Nil(err)
}

func TestMapGeneratesError(t *testing.T) {
	errorValue := "Do"
	expectedError := errors.New("this is an example error")

	inputSlice := []string{"DogCows", "goes", "Moof!", errorValue, "you", "like", "Clarus", "the", "DogCow?"}

	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	mapFunction := func(value string) (int, error) {
		if value == errorValue {
			return 0, expectedError
		}

		return len(value), nil
	}

	_, err := link.Map(mapFunction).Slice()

	assert.Equal(t, expectedError, err)
}

func TestFilter(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	expectedSlice := []int{7, 9, 6, 8}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	filterFunction := func(value int) (bool, error) {
		return value > 5, nil
	}

	actualSlice, err := link.Filter(filterFunction).Slice()

	assert.Equal(expectedSlice, actualSlice)
	assert.Nil(err)
}

func TestFilterHasError(t *testing.T) {
	assert := assert.New(t)

	errorValue := 9
	expectedError := errors.New("an example error")
	inputSlice := []int{7, 4, 2, 3, errorValue, 5, 6, 0, 8, 1}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	filterFunction := func(value int) (bool, error) {
		if value == errorValue {
			return false, expectedError
		}
		return value > 5, nil
	}

	_, err := link.Filter(filterFunction).Slice()

	assert.Equal(expectedError, err)
}

func TestSkip(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	howManyToSkip := 3
	actualSlice, err := link.Skip(howManyToSkip).Slice()

	assert.Equal(inputSlice[howManyToSkip:], actualSlice)
	assert.Nil(err)
}

func TestSkipLargerThanSlice(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	actualSlice, err := link.Skip(len(inputSlice) + 1).Slice()

	assert.Equal([]int{}, actualSlice)
	assert.Nil(err)
}

func TestLimit(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	howManyToKeep := 6
	actualSlice, err := link.Limit(howManyToKeep).Slice()

	assert.Equal(inputSlice[:howManyToKeep], actualSlice)
	assert.Nil(err)
}

func TestLimitLargerThanSlice(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	actualSlice, err := link.Limit(len(inputSlice) + 1).Slice()

	assert.Equal(inputSlice, actualSlice)
	assert.Nil(err)
}

func TestDistinctFunc(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 7, 3, 7, 9, 5, 5, 2, 6, 0, 8, 1}
	expectedSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	actualSlice, err := link.DistinctFunc(func(value int) int { return value }).Slice()

	assert.Equal(expectedSlice, actualSlice)
	assert.Nil(err)
}

func TestFlattenWithSliceOfSlice(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []any{[]int{1, 2, 3}, []int{4, 5, 6}, []int{7, 8, 9}}
	expectedSlice := []any{1, 2, 3, 4, 5, 6, 7, 8, 9}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	actualSlice, err := link.Flatten[any]().Slice()

	assert.Equal(expectedSlice, actualSlice)
	assert.Nil(err)
}

func TestFlattenWithSliceMix(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []any{[]int{1, 2, 3}, 4, []int{7, 8, 9}}
	expectedSlice := []any{1, 2, 3, 4, 7, 8, 9}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	actualSlice, err := link.Flatten[any]().Slice()

	assert.Equal(expectedSlice, actualSlice)
	assert.Nil(err)
}

func TestFlattenWithArray(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []any{[...]int{1, 2, 3}, 4, [...]int{7, 8, 9}}
	expectedSlice := []any{1, 2, 3, 4, 7, 8, 9}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	actualSlice, err := link.Flatten[any]().Slice()

	assert.Equal(expectedSlice, actualSlice)
	assert.Nil(err)
}

func TestFlattenWithChannel(t *testing.T) {
	assert := assert.New(t)

	firstChannel := createTestIntChannel([]int{1, 2, 3})
	secondChannel := createTestIntChannel([]int{7, 8, 9})

	inputSlice := []any{firstChannel, 4, secondChannel}
	expectedSlice := []any{1, 2, 3, 4, 7, 8, 9}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	actualSlice, err := link.Flatten[any]().Slice()

	assert.Equal(expectedSlice, actualSlice)
	assert.Nil(err)
}

func TestFlattenWithSliceAndMap(t *testing.T) {
	assert := assert.New(t)

	key1 := 4
	value1 := 5
	key2 := 6
	value2 := 7

	innerMap := map[int]int{
		key1: value1,
		key2: value2,
	}
	inputSlice := []any{[]int{1, 2, 3}, innerMap, []int{7, 8, 9}}

	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	actualSlice, err := link.Flatten[any]().Slice()

	assert.Nil(err)
	assert.Len(actualSlice, 8)
	// First three are 1, 2, 3
	assert.Equal(1, actualSlice[0])
	assert.Equal(2, actualSlice[1])
	assert.Equal(3, actualSlice[2])
	// Last three are 7, 8, 9
	assert.Equal(7, actualSlice[5])
	assert.Equal(8, actualSlice[6])
	assert.Equal(9, actualSlice[7])
	// Middle two are the map entries (order not guaranteed)
	seenKeys := map[int]int{}
	for _, entry := range actualSlice[3:5] {
		kv, ok := entry.(keyvalue.KeyValuer[any, any])
		assert.True(ok)
		seenKeys[kv.Key().(int)] = kv.Value().(int)
	}
	assert.Equal(value1, seenKeys[key1])
	assert.Equal(value2, seenKeys[key2])
}

func TestFlattenTypedInts(t *testing.T) {
	assert := assert.New(t)

	inputSlice := [][]int{{1, 2, 3}, {4, 5, 6}}
	expectedSlice := []int{1, 2, 3, 4, 5, 6}

	actualSlice, err := FromSlice(inputSlice).Flatten[int]().Slice()

	assert.Equal(expectedSlice, actualSlice)
	assert.Nil(err)
}

func TestFlattenTypedMismatchInjectsError(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []any{[]int{1, 2}, "oops", []int{3}}

	actualSlice, err := FromSlice(inputSlice).Flatten[int]().Slice()

	assert.Equal([]int{1, 2}, actualSlice)
	assert.NotNil(err)
	assert.Contains(err.Error(), "not assignable to int")
}

func TestSort(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	expectedSlice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	returnLessFunction := func(sliceToSort []int) func(int, int) bool {
		return func(i int, j int) bool {
			return sliceToSort[i] < sliceToSort[j]
		}
	}
	actualSlice, err := link.Sort(returnLessFunction).Slice()

	assert.Equal(expectedSlice, actualSlice)
	assert.Nil(err)
}

func TestSortHasError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("this is an example error")
	errorValue := 9
	inputSlice := []int{7, 4, 2, 3, errorValue, 5, 6, 0, 8, 1}
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := newLink(generation)

	returnLessFunction := func(sliceToSort []int) func(int, int) bool {
		return func(i int, j int) bool {
			return sliceToSort[i] < sliceToSort[j]
		}
	}
	_, err := link.Sort(returnLessFunction).Slice()

	assert.Equal(expectedError, err)
}

func TestReverse(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	expectedSlice := []int{1, 8, 0, 6, 5, 9, 3, 2, 4, 7}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	actualSlice, err := link.Reverse().Slice()

	assert.Equal(expectedSlice, actualSlice)
	assert.Nil(err)
}

func TestReverseHasError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("this is an example error")
	errorValue := 9
	inputSlice := []int{7, 4, 2, 3, errorValue, 5, 6, 0, 8, 1}
	generation := createGeneratorWithError(inputSlice, errorValue, expectedError)
	link := newLink(generation)

	_, err := link.Reverse().Slice()

	assert.Equal(expectedError, err)
}

func createTestIntChannel(intSlice []int) chan int {
	intChannel := make(chan int)

	go func() {
		for _, currentInt := range intSlice {
			intChannel <- currentInt
		}
		close(intChannel)
	}()

	return intChannel
}

func createGeneratorWithError(intSlice []int, valueToErrorOn int, errorToReturn error) func() (int, error) {
	generation := generator.FromSlice(intSlice)

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
