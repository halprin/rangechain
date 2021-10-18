package generator

import (
	"errors"
	"github.com/halprin/slice-chain/helper"
	"reflect"
)

var Exhausted = errors.New("generator exhausted")

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

func FromChannel(channel chan interface{}) func() (interface{}, error) {
	if !helper.IsChannel(channel) {
		panic("non-channel type provided")
	}

	return func() (interface{}, error) {
		for value := range channel {
			return value, nil
		}

		return 0, Exhausted
	}
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
