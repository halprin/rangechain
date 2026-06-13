package rangechain

import (
	"errors"
	"fmt"
	"github.com/halprin/rangechain/internal/generator"
	"reflect"
)

// Flatten will iterate over all the values in the chain, but any value encountered that is a range-able container itself will also have its values iterated over first before continuing with the remaining values in the chain.  Maps flatten to its `keyvalue.KeyValuer[any, any]` key and value pairs.
//
// The caller specifies the output element type `U`.  Each emitted inner value (or scalar) is type-asserted to `U`; on mismatch, an error is injected into the chain at that point and downstream terminations propagate it.
func (receiver *Link[T]) Flatten[U any]() *Link[U] {
	var currentGenerator func() (any, error)

	flattenGenerator := func() (U, error) {
		var innerValue any
		var err error

		for innerValue == nil {
			if currentGenerator == nil {
				var currentValue T
				currentValue, err = receiver.generator()
				if err != nil {
					var zero U
					return zero, err
				}

				currentAny := any(currentValue)
				if isSlice(currentAny) || isArray(currentAny) {
					currentGenerator = sliceOrArrayAnyGenerator(currentAny)
				} else if isChannel(currentAny) {
					currentGenerator = channelAnyGenerator(currentAny)
				} else if isMap(currentAny) {
					currentGenerator = mapAnyGenerator(currentAny)
				} else {
					innerValue = currentAny
					break
				}
			}

			innerValue, err = currentGenerator()
			if errors.Is(err, generator.Exhausted) {
				innerValue = nil
				currentGenerator = nil
			}
		}

		return assertTo[U](innerValue)
	}

	return newLink(flattenGenerator)
}

func assertTo[U any](v any) (U, error) {
	if casted, ok := v.(U); ok {
		return casted, nil
	}
	var zero U
	return zero, fmt.Errorf("flatten: element %v (type %T) is not assignable to %v", v, v, reflect.TypeFor[U]())
}

func isSlice(value any) bool {
	return reflect.ValueOf(value).Kind() == reflect.Slice
}

func isArray(value any) bool {
	return reflect.ValueOf(value).Kind() == reflect.Array
}

func isChannel(value any) bool {
	return reflect.ValueOf(value).Kind() == reflect.Chan
}

func isMap(value any) bool {
	return reflect.ValueOf(value).Kind() == reflect.Map
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
