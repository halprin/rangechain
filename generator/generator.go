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

	concreteValue := reflect.ValueOf(slice)

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
