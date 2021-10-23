package intermediate

import "github.com/halprin/rangechain/generator"

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

func (receiver *Link) FilterParallel(filterFunction func(interface{}) bool) *Link {
	computedValues := false
	var mappedReturnValues []chan interface{}
	currentIndex := 0

	filterGenerator := func() (interface{}, error) {
		if !computedValues {
			//run the map function against all the items in the generator on the first call to this generator
			mappedReturnValues = filterFunctionAgainstEntireGenerator(receiver.generator, filterFunction)
			computedValues = true
		}

		//go through the return values until you find an item that stays
		for {
			if currentIndex >= len(mappedReturnValues) {
				return 0, generator.Exhausted
			}

			value := <- mappedReturnValues[currentIndex]
			valueStays := (<- mappedReturnValues[currentIndex]).(bool)
			currentIndex++

			if valueStays {
				return value, nil
			}
		}
	}

	return NewLink(filterGenerator)
}

func filterFunctionAgainstEntireGenerator(generatorToParallelize func() (interface{}, error), filterFunction func(interface{}) bool) []chan interface{} {
	var mappedReturnValues []chan interface{}

	for {
		valueToFilter, err := generatorToParallelize()
		if err != nil {
			break
		}

		mappedReturnValue := make(chan interface{})
		mappedReturnValues = append(mappedReturnValues, mappedReturnValue)

		go pipeInputValueAndReturnValueToChannel(filterFunction, valueToFilter, mappedReturnValue)
	}

	return mappedReturnValues
}

func pipeInputValueAndReturnValueToChannel(mapFunction func(interface{}) bool, valueToMap interface{}, returnValueChannel chan interface{}) {
	returnValueChannel <- valueToMap
	returnValueChannel <- mapFunction(valueToMap)
	close(returnValueChannel)
}
