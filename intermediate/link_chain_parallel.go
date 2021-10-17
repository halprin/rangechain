package intermediate

import "github.com/halprin/slice-chain/generator"

func (receiver *Link) MapParallel(mapFunction func(interface{}) interface{}) *Link {
	computedValues := false
	var mappedReturnValues []chan interface{}
	currentIndex := 0

	mapGenerator := func() (interface{}, error) {
		if !computedValues {
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
