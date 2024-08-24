package rangechain

import (
	"errors"
	"github.com/halprin/rangechain/internal/generator"
	"github.com/halprin/rangechain/internal/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	var expectedOutput []interface{}
	for _, stringValue := range inputSlice {
		expectedOutput = append(expectedOutput, len(stringValue))
	}

	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	mapFunction := func(value interface{}) (interface{}, error) {
		stringValue := value.(string)
		return len(stringValue), nil
	}

	actualSlice, err := link.Map(mapFunction).Slice()

	assert.Equal(expectedOutput, actualSlice)
	assert.Nil(err)
}

func TestMapGeneratesError(t *testing.T) {
	errorValue := "Do"
	expectedError := errors.New("this is an example error")

	inputSlice := []string{"DogCows", "goes", "Moof!", errorValue, "you", "like", "Clarus", "the", "DogCow?"}
	var expectedOutput []interface{}
	for _, stringValue := range inputSlice {
		expectedOutput = append(expectedOutput, len(stringValue))
	}

	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	mapFunction := func(value interface{}) (interface{}, error) {
		stringValue := value.(string)
		if stringValue == errorValue {
			return 0, expectedError
		}

		return len(stringValue), nil
	}

	_, err := link.Map(mapFunction).Slice()

	assert.Equal(t, expectedError, err)
}

func TestFilter(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	expectedSlice := helper.InterfaceSlice([]int{7, 9, 6, 8})
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	filterFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		return intValue > 5, nil
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

	filterFunction := func(value interface{}) (bool, error) {
		intValue := value.(int)
		if intValue == errorValue {
			return false, expectedError
		}
		return intValue > 5, nil
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

	assert.Equal(helper.InterfaceSlice(inputSlice[howManyToSkip:]), actualSlice)
	assert.Nil(err)
}

func TestSkipLargerThanSlice(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	actualSlice, err := link.Skip(len(inputSlice) + 1).Slice()

	assert.Equal([]interface{}{}, actualSlice)
	assert.Nil(err)
}

func TestLimit(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	howManyToKeep := 6
	actualSlice, err := link.Limit(howManyToKeep).Slice()

	assert.Equal(helper.InterfaceSlice(inputSlice[:howManyToKeep]), actualSlice)
	assert.Nil(err)
}

func TestLimitLargerThanSlice(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	actualSlice, err := link.Limit(len(inputSlice) + 1).Slice()

	assert.Equal(helper.InterfaceSlice(inputSlice), actualSlice)
	assert.Nil(err)
}

func TestDistinct(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 7, 3, 7, 9, 5, 5, 2, 6, 0, 8, 1}
	expectedSlice := helper.InterfaceSlice([]int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1})
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	actualSlice, err := link.Distinct().Slice()

	assert.Equal(expectedSlice, actualSlice)
	assert.Nil(err)
}

//func TestFlattenWithSliceOfSlice(t *testing.T) {
//	assert := assert.New(t)
//
//	inputSlice := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
//	expectedSlice := helper.InterfaceSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
//	generation := generator.FromSlice(inputSlice)
//	link := newLink(generation)
//
//	actualSlice, err := link.Flatten().Slice()
//
//	assert.Equal(expectedSlice, actualSlice)
//	assert.Nil(err)
//}

//func TestFlattenWithSliceMix(t *testing.T) {
//	assert := assert.New(t)
//
//	inputSlice := []interface{}{[]int{1, 2, 3}, 4, []int{7, 8, 9}}
//	expectedSlice := helper.InterfaceSlice([]int{1, 2, 3, 4, 7, 8, 9})
//	generation := generator.FromSlice(inputSlice)
//	link := newLink(generation)
//
//	actualSlice, err := link.Flatten().Slice()
//
//	assert.Equal(expectedSlice, actualSlice)
//	assert.Nil(err)
//}

//func TestFlattenWithArray(t *testing.T) {
//	assert := assert.New(t)
//
//	inputSlice := []interface{}{[...]int{1, 2, 3}, 4, [...]int{7, 8, 9}}
//	expectedSlice := helper.InterfaceSlice([]int{1, 2, 3, 4, 7, 8, 9})
//	generation := generator.FromSlice(inputSlice)
//	link := newLink(generation)
//
//	actualSlice, err := link.Flatten().Slice()
//
//	assert.Equal(expectedSlice, actualSlice)
//	assert.Nil(err)
//}

//func TestFlattenWithChannel(t *testing.T) {
//	assert := assert.New(t)
//
//	firstChannel := createTestIntChannel([]int{1, 2, 3})
//	secondChannel := createTestIntChannel([]int{7, 8, 9})
//
//	inputSlice := []interface{}{firstChannel, 4, secondChannel}
//	expectedSlice := helper.InterfaceSlice([]int{1, 2, 3, 4, 7, 8, 9})
//	generation := generator.FromSlice(inputSlice)
//	link := newLink(generation)
//
//	actualSlice, err := link.Flatten().Slice()
//
//	assert.Equal(expectedSlice, actualSlice)
//	assert.Nil(err)
//}

//func TestFlattenWithSliceAndMap(t *testing.T) {
//	assert := assert.New(t)
//
//	key1 := 4
//	value1 := 5
//	key2 := 6
//	value2 := 7
//
//	innerMap := map[int]int{
//		key1: value1,
//		key2: value2,
//	}
//	inputSlice := []interface{}{[]int{1, 2, 3}, innerMap, []int{7, 8, 9}}
//	expectedSlice := []interface{}{
//		1,
//		2,
//		3,
//		&testKeyValue{
//			TheKey:   key1,
//			TheValue: value1,
//		},
//		&testKeyValue{
//			TheKey:   key2,
//			TheValue: value2,
//		},
//		7,
//		8,
//		9,
//	}
//
//	generation := generator.FromSlice(inputSlice)
//	link := newLink(generation)
//
//	actualSlice, err := link.Flatten().Slice()
//
//	//not testing the order of the entire expected slice because we are not guaranteed the order in which a map is iterated over
//	assertEqualsBasedOnKeyValuerInterface(t, expectedSlice, actualSlice)
//	assert.Nil(err)
//	//test the order for the non-map flattened parts
//	assert.Equal(expectedSlice[0], actualSlice[0])
//	assert.Equal(expectedSlice[1], actualSlice[1])
//	assert.Equal(expectedSlice[2], actualSlice[2])
//	assert.Equal(expectedSlice[5], actualSlice[5])
//	assert.Equal(expectedSlice[6], actualSlice[6])
//	assert.Equal(expectedSlice[7], actualSlice[7])
//}

func TestSort(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	expectedSlice := helper.InterfaceSlice([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	generation := generator.FromSlice(inputSlice)
	link := newLink(generation)

	returnLessFunction := func(sliceToSort []interface{}) func(int, int) bool {
		return func(i int, j int) bool {
			iItem := sliceToSort[i].(int)
			jItem := sliceToSort[j].(int)
			return iItem < jItem
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

	returnLessFunction := func(sliceToSort []interface{}) func(int, int) bool {
		return func(i int, j int) bool {
			iItem := sliceToSort[i].(int)
			jItem := sliceToSort[j].(int)
			return iItem < jItem
		}
	}
	_, err := link.Sort(returnLessFunction).Slice()

	assert.Equal(expectedError, err)
}

func TestReverse(t *testing.T) {
	assert := assert.New(t)

	inputSlice := []int{7, 4, 2, 3, 9, 5, 6, 0, 8, 1}
	expectedSlice := helper.InterfaceSlice([]int{1, 8, 0, 6, 5, 9, 3, 2, 4, 7})
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

func createTestIntChannel(intSlice []int) chan interface{} {
	intChannel := make(chan interface{})

	go func() {
		for _, currentInt := range intSlice {
			intChannel <- currentInt
		}
		close(intChannel)
	}()

	return intChannel
}

func createGeneratorWithError(intSlice []int, valueToErrorOn int, errorToReturn error) func() (interface{}, error) {
	generation := generator.FromSlice(intSlice)

	return func() (interface{}, error) {
		value, err := generation()

		if err != nil {
			return value, err
		} else if value == valueToErrorOn {
			return 0, errorToReturn
		}

		return value, err
	}
}
