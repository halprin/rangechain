package intermediate

import "github.com/halprin/rangechain/internal/generator"

// MapParallel will run the `mapFunction` parameter function against all the values in the chain in parallel.  In that function, return what you want to change the value into or an optional error if an error is encountered.  There is overhead to running in parallel so benchmark to ensure you benefit from this version.
func (receiver *Link) MapParallel(mapFunction func(interface{}) (interface{}, error)) *Link {
	computedValues := false
	var mappedReturnValues []chan interface{}
	var mappedErrorValues []chan error
	currentIndex := 0

	mapGenerator := func() (interface{}, error) {
		if !computedValues {
			//run the map function against all the items in the generator on the first call to this generator
			mappedReturnValues, mappedErrorValues = mapFunctionAgainstEntireGenerator(receiver.generator, mapFunction)
			computedValues = true
		}

		if currentIndex >= len(mappedReturnValues) {
			return 0, generator.Exhausted
		}

		value := <- mappedReturnValues[currentIndex]
		err := <- mappedErrorValues[currentIndex]
		currentIndex++

		return value, err
	}

	return NewLink(mapGenerator)
}

func mapFunctionAgainstEntireGenerator(generatorToParallelize func() (interface{}, error), mapFunction func(interface{}) (interface{}, error)) ([]chan interface{}, []chan error) {
	var mappedReturnValues []chan interface{}
	var mappedErrorValues []chan error

	for {
		valueToMap, err := generatorToParallelize()
		if err != nil {
			break
		}

		mappedReturnValue := make(chan interface{})
		mappedReturnValues = append(mappedReturnValues, mappedReturnValue)
		mappedErrorValue := make(chan error)
		mappedErrorValues = append(mappedErrorValues, mappedErrorValue)

		go pipeReturnAndErrorValueToChannels(mapFunction, valueToMap, mappedReturnValue, mappedErrorValue)
	}

	return mappedReturnValues, mappedErrorValues
}

func pipeReturnAndErrorValueToChannels(mapFunction func(interface{}) (interface{}, error), valueToMap interface{}, returnValueChannel chan interface{}, returnErrorChannel chan error) {
	mappedValue, err := mapFunction(valueToMap)
	returnValueChannel <- mappedValue
	returnErrorChannel <- err
	close(returnValueChannel)
	close(returnErrorChannel)
}

// FilterParallel will run the `filterFunction` parameter function against all the values in the chain in parallel.  In that function, on return of true, the value will stay in the chain, or on false, the value will be removed from the chain.  There is overhead to running in parallel so benchmark to ensure you benefit from this version.
func (receiver *Link) FilterParallel(filterFunction func(interface{}) (bool, error)) *Link {
	computedValues := false
	var mappedReturnValues []chan interface{}
	var mappedErrorValues []chan error
	currentIndex := 0

	filterGenerator := func() (interface{}, error) {
		if !computedValues {
			//run the map function against all the items in the generator on the first call to this generator
			mappedReturnValues, mappedErrorValues = filterFunctionAgainstEntireGenerator(receiver.generator, filterFunction)
			computedValues = true
		}

		//go through the return values until you find an item that stays
		for {
			if currentIndex >= len(mappedReturnValues) {
				return 0, generator.Exhausted
			}

			value := <- mappedReturnValues[currentIndex]
			valueStays := (<- mappedReturnValues[currentIndex]).(bool)
			err := <- mappedErrorValues[currentIndex]
			currentIndex++

			if err != nil {
				return value, err
			} else if valueStays {
				return value, nil
			}
		}
	}

	return NewLink(filterGenerator)
}

func filterFunctionAgainstEntireGenerator(generatorToParallelize func() (interface{}, error), filterFunction func(interface{}) (bool, error)) ([]chan interface{}, []chan error) {
	var mappedReturnValues []chan interface{}
	var mappedErrorValues []chan error

	for {
		valueToFilter, err := generatorToParallelize()
		if err != nil {
			break
		}

		mappedReturnValue := make(chan interface{})
		mappedReturnValues = append(mappedReturnValues, mappedReturnValue)
		mappedErrorValue := make(chan error)
		mappedErrorValues = append(mappedErrorValues, mappedErrorValue)

		go pipeInputValueAndReturnValueToChannel(filterFunction, valueToFilter, mappedReturnValue, mappedErrorValue)
	}

	return mappedReturnValues, mappedErrorValues
}

func pipeInputValueAndReturnValueToChannel(filterFunction func(interface{}) (bool, error), valueToMap interface{}, returnValueChannel chan interface{}, returnErrorChannel chan error) {
	filtered, err := filterFunction(valueToMap)
	returnValueChannel <- valueToMap
	returnValueChannel <- filtered
	returnErrorChannel <- err
	close(returnValueChannel)
	close(returnErrorChannel)
}
