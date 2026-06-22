package rangechain

import "github.com/halprin/rangechain/internal/generator"

// MapParallel is like Map, but invocations run concurrently. There is overhead to running in parallel so benchmark to ensure you benefit from this version.
func (receiver *Link[T]) MapParallel[U any](mapFunction func(T) (U, error)) *Link[U] {
	computedValues := false
	var mappedReturnValues []chan U
	var mappedErrorValues []chan error
	currentIndex := 0

	mapGenerator := func() (U, error) {
		if !computedValues {
			mappedReturnValues, mappedErrorValues = mapFunctionAgainstEntireGenerator(receiver.generator, mapFunction)
			computedValues = true
		}

		if currentIndex >= len(mappedReturnValues) {
			var zero U
			return zero, generator.Exhausted
		}

		value := <-mappedReturnValues[currentIndex]
		err := <-mappedErrorValues[currentIndex]
		currentIndex++

		return value, err
	}

	return newLink(mapGenerator)
}

func mapFunctionAgainstEntireGenerator[T, U any](generatorToParallelize func() (T, error), mapFunction func(T) (U, error)) ([]chan U, []chan error) {
	var mappedReturnValues []chan U
	var mappedErrorValues []chan error

	for {
		valueToMap, err := generatorToParallelize()
		if err != nil {
			break
		}

		mappedReturnValue := make(chan U)
		mappedReturnValues = append(mappedReturnValues, mappedReturnValue)
		mappedErrorValue := make(chan error)
		mappedErrorValues = append(mappedErrorValues, mappedErrorValue)

		go pipeReturnAndErrorValueToChannels(mapFunction, valueToMap, mappedReturnValue, mappedErrorValue)
	}

	return mappedReturnValues, mappedErrorValues
}

func pipeReturnAndErrorValueToChannels[T, U any](mapFunction func(T) (U, error), valueToMap T, returnValueChannel chan U, returnErrorChannel chan error) {
	mappedValue, err := mapFunction(valueToMap)
	returnValueChannel <- mappedValue
	returnErrorChannel <- err
	close(returnValueChannel)
	close(returnErrorChannel)
}

type filterResult[T any] struct {
	value T
	keep  bool
}

// FilterParallel is like Filter, but invocations run concurrently. There is overhead to running in parallel so benchmark to ensure you benefit from this version.
func (receiver *Link[T]) FilterParallel(filterFunction func(T) (bool, error)) *Link[T] {
	computedValues := false
	var resultChannels []chan filterResult[T]
	var errorChannels []chan error
	currentIndex := 0

	filterGenerator := func() (T, error) {
		if !computedValues {
			resultChannels, errorChannels = filterFunctionAgainstEntireGenerator(receiver.generator, filterFunction)
			computedValues = true
		}

		for {
			if currentIndex >= len(resultChannels) {
				var zero T
				return zero, generator.Exhausted
			}

			result := <-resultChannels[currentIndex]
			err := <-errorChannels[currentIndex]
			currentIndex++

			if err != nil {
				return result.value, err
			} else if result.keep {
				return result.value, nil
			}
		}
	}

	return newLink(filterGenerator)
}

func filterFunctionAgainstEntireGenerator[T any](generatorToParallelize func() (T, error), filterFunction func(T) (bool, error)) ([]chan filterResult[T], []chan error) {
	var resultChannels []chan filterResult[T]
	var errorChannels []chan error

	for {
		valueToFilter, err := generatorToParallelize()
		if err != nil {
			break
		}

		resultChannel := make(chan filterResult[T])
		resultChannels = append(resultChannels, resultChannel)
		errorChannel := make(chan error)
		errorChannels = append(errorChannels, errorChannel)

		go pipeFilterResultToChannel(filterFunction, valueToFilter, resultChannel, errorChannel)
	}

	return resultChannels, errorChannels
}

func pipeFilterResultToChannel[T any](filterFunction func(T) (bool, error), valueToFilter T, resultChannel chan filterResult[T], errorChannel chan error) {
	keep, err := filterFunction(valueToFilter)
	resultChannel <- filterResult[T]{value: valueToFilter, keep: keep}
	errorChannel <- err
	close(resultChannel)
	close(errorChannel)
}
