package intermediate

import "github.com/halprin/rangechain/generator"

func (receiver *Link) MapParallel(mapFunction func(interface{}) interface{}) *Link {
	computedValues := false
	var mappedReturnValues []chan interface{}
	currentIndex := 0

	mapGenerator := func() (interface{}, error) {
		if !computedValues {
			//run the map function against all the items in the generator on the first call to this generator
			mappedReturnValues = mapFunctionAgainstEntireGenerator(receiver.generator, mapFunction)
			computedValues = true
		}

		if currentIndex >= len(mappedReturnValues) {
			return 0, generator.Exhausted
		}

		value := <- mappedReturnValues[currentIndex]
		currentIndex++

		return value, nil
	}

	return NewLink(mapGenerator)
}

func mapFunctionAgainstEntireGenerator(generatorToParallelize func() (interface{}, error), mapFunction func(interface{}) interface{}) []chan interface{} {
	var mappedReturnValues []chan interface{}

	for {
		valueToMap, err := generatorToParallelize()
		if err != nil {
			break
		}

		mappedReturnValue := make(chan interface{})
		mappedReturnValues = append(mappedReturnValues, mappedReturnValue)

		go pipeReturnValueToChannel(mapFunction, valueToMap, mappedReturnValue)
	}

	return mappedReturnValues
}

func pipeReturnValueToChannel(mapFunction func(interface{}) interface{}, valueToMap interface{}, returnValueChannel chan interface{}) {
	returnValueChannel <- mapFunction(valueToMap)
	close(returnValueChannel)
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
