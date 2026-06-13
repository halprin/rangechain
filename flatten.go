package rangechain

import (
	"errors"
	"github.com/halprin/rangechain/internal/generator"
	"github.com/halprin/rangechain/internal/helper"
	"reflect"
)

// Flatten will iterate over all the values in the chain, but any value encountered that is a range-able container itself will also have its values iterated over first before continuing with the remaining values in the chain.  Maps flatten to its `keyvalue.KeyValuer` key and value pairs.
//
// Because Flatten requires the chain to hold heterogeneous container types, it is only available on `*Link[any]`.
func Flatten(receiver *Link[any]) *Link[any] {
	var currentGenerator func() (any, error)

	flattenGenerator := func() (any, error) {
		var innerValue any
		var err error

		for innerValue == nil {
			if currentGenerator == nil {
				var currentValue any
				currentValue, err = receiver.generator()
				if err != nil {
					return nil, err
				}

				if helper.IsSlice(currentValue) || helper.IsArray(currentValue) {
					currentGenerator = sliceOrArrayAnyGenerator(currentValue)
				} else if helper.IsChannel(currentValue) {
					currentGenerator = channelAnyGenerator(currentValue)
				} else if helper.IsMap(currentValue) {
					currentGenerator = mapAnyGenerator(currentValue)
				} else {
					innerValue = currentValue
					break
				}
			}

			innerValue, err = currentGenerator()
			if errors.Is(err, generator.Exhausted) {
				innerValue = nil
				currentGenerator = nil
			}
		}

		return innerValue, err
	}

	return newLink(flattenGenerator)
}

// sliceOrArrayAnyGenerator builds an any-typed generator over a slice or array of unknown element type, using reflection.
func sliceOrArrayAnyGenerator(sliceOrArray any) func() (any, error) {
	concreteValue := reflect.ValueOf(sliceOrArray)
	currentIndex := 0

	return func() (any, error) {
		if currentIndex >= concreteValue.Len() {
			return nil, generator.Exhausted
		}

		value := concreteValue.Index(currentIndex).Interface()
		currentIndex++

		return value, nil
	}
}

// channelAnyGenerator builds an any-typed generator over a channel of unknown element type, using reflection.
func channelAnyGenerator(channel any) func() (any, error) {
	concreteValue := reflect.ValueOf(channel)

	return func() (any, error) {
		value, ok := concreteValue.Recv()
		if !ok {
			return nil, generator.Exhausted
		}

		return value.Interface(), nil
	}
}

// mapAnyGenerator builds an any-typed generator over a map of unknown key/value types, using reflection.  Each emitted value is a `flattenMapTuple` that implements `keyvalue.KeyValuer[any, any]`.
func mapAnyGenerator(aMap any) func() (any, error) {
	concreteValue := reflect.ValueOf(aMap)
	mapIterator := concreteValue.MapRange()

	return func() (any, error) {
		hasNext := mapIterator.Next()
		if !hasNext {
			return nil, generator.Exhausted
		}

		return &flattenMapTuple{
			TheKey:   mapIterator.Key().Interface(),
			TheValue: mapIterator.Value().Interface(),
		}, nil
	}
}

type flattenMapTuple struct {
	TheKey   any
	TheValue any
}

func (m *flattenMapTuple) Key() any {
	return m.TheKey
}

func (m *flattenMapTuple) Value() any {
	return m.TheValue
}
