// Package generator has many helper methods for generating values in a chain.
package generator

import (
	"errors"
	"github.com/halprin/rangechain/internal/helper"
	"reflect"
)

// Exhausted is returned as an expected error from the generators to designate an end of the generator.
var Exhausted = errors.New("generator exhausted")

// MapTuple is used to represent the key and value pairs of a map when iterating over a map.
type MapTuple struct {
	Key   interface{}
	Value interface{}
}

// FromSlice creates a generator for a slice.
func FromSlice(slice interface{}) func() (interface{}, error) {
	if !helper.IsSlice(slice) {
		panic("non-slice type provided")
	}

	return generatorFromSliceOrArray(slice)
}

// FromArray creates a generator for an array.
func FromArray(array interface{}) func() (interface{}, error) {
	if !helper.IsArray(array) {
		panic("non-array type provided")
	}

	return generatorFromSliceOrArray(array)
}

// FromChannel creates a generator for a channel.
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

// FromMap creates a generator for a map.
func FromMap(aMap interface{}) func() (interface{}, error) {
	if !helper.IsMap(aMap) {
		panic("non-map type provided")
	}

	concreteValue := reflect.ValueOf(aMap)
	mapIterator := concreteValue.MapRange()

	return func() (interface{}, error) {
		hasNext := mapIterator.Next()
		if !hasNext {
			return 0, Exhausted
		}

		mapTuple := MapTuple{
			Key:   mapIterator.Key().Interface(),
			Value: mapIterator.Value().Interface(),
		}

		return mapTuple, nil
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
