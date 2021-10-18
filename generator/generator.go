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
