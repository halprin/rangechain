package generator

import (
	"errors"
	"github.com/halprin/slice-chain/helper"
	"reflect"
)

var Exhausted = errors.New("generator exhausted")

type MapTuple struct {
	Key   interface{}
	Value interface{}
}

func FromSlice(slice interface{}) func() (interface{}, error) {
	if !helper.IsSlice(slice) {
		panic("non-slice type provided")
	}

	return generatorFromSliceOrArray(slice)
}

func FromArray(array interface{}) func() (interface{}, error) {
	if !helper.IsArray(array) {
		panic("non-array type provided")
	}

	return generatorFromSliceOrArray(array)
}

func FromChannel(channel interface{}) func() (interface{}, error) {
	if !helper.IsChannel(channel) {
		panic("non-channel type provided")
	}

	concreteValue := reflect.ValueOf(channel)

	return func() (interface{}, error) {
		for {
			value, ok := concreteValue.Recv()
			if !ok {
				return 0, Exhausted
			}

			return value.Interface(), nil
		}
	}
}

func FromMap(aMap interface{}) func() (interface{}, error) {
	if !helper.IsMap(aMap) {
		panic("non-map type provided")
	}

	actualMap := aMap.(map[interface{}]interface{})
	mapChannel := make(chan MapTuple)

	//convert to using channels because...
	//there is no way to directly access into a map using an index,
	//nor a way to have the "for range" start where the last call to the generator left off.

	go func() {
		for key, value := range actualMap {
			mapTuple := MapTuple{
				Key:   key,
				Value: value,
			}

			mapChannel <- mapTuple
		}
	}()

	return FromChannel(mapChannel)
}

func generatorFromSliceOrArray(sliceOrArray interface{}) func() (interface{}, error) {
	concreteValue := reflect.ValueOf(sliceOrArray)

	currentIndex := 0

	return func() (interface{}, error) {
		if currentIndex >= concreteValue.Len() {
			return 0, Exhausted
		}

		value := concreteValue.Index(currentIndex).Interface()
		currentIndex++

		return value, nil
	}
}
