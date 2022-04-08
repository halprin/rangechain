package rangechain

import (
	"errors"
	"github.com/halprin/rangechain/internal/generator"
	"github.com/halprin/rangechain/internal/helper"
	"sort"
)

// Map will run the `mapFunction` parameter function parameter against all the values in the chain.  In that function, return what you want to change the value into or an optional error if an error is encountered.
func (receiver *Link[T]) Map(mapFunction func(T) (interface{}, error)) *Link {
	mapGenerator := func() (interface{}, error) {
		valueToMap, err := receiver.generator()
		if err != nil {
			return 0, err
		}

		return mapFunction(valueToMap)
	}

	return newLink(mapGenerator)
}

func MapFunction[T any, S any](link *Link[T], mapFunction func(T) (S, error)) *Link[S] {
	mapGenerator := func() (S, error) {
		valueToMap, err := link.generator()
		if err != nil {
			var zeroReturnValue S
			return zeroReturnValue, err
		}

		return mapFunction(valueToMap)
	}

	return newLink(mapGenerator)
}

// Filter will run the `filterFunction` parameter function against all the values in the chain.  In that function, on return of true, the value will stay, or on false, the value will be removed.
func (receiver *Link[T]) Filter(filterFunction func(T) (bool, error)) *Link[T] {
	filterGenerator := func() (T, error) {
		//go through the generator until you find an item that stays
		for {
			valueToFilter, err := receiver.generator()
			if err != nil {
				var zeroReturnValue T
				return zeroReturnValue, err
			}

			valueStays, err := filterFunction(valueToFilter)

			if err != nil {
				return valueToFilter, err
			} else if valueStays {
				return valueToFilter, nil
			}
		}
	}

	return newLink(filterGenerator)
}

// Skip skips over the parameter `skipNumber` number of values and effectively removes them from the chain.  Also skips over any errors previously generated.
func (receiver *Link[T]) Skip(skipNumber int) *Link[T] {
	for count := 0; count < skipNumber; count++ {
		_, _ = receiver.generator()
	}

	return newLink(receiver.generator)
}

// Limit stops the chain after the parameter `keepSize` number of values.  Any elements afterward are effectively removed.
func (receiver *Link[T]) Limit(keepSize int) *Link[T] {
	itemsSeen := 0

	limitGenerator := func() (T, error) {
		if itemsSeen >= keepSize {
			var zeroReturnValue T
			return zeroReturnValue, generator.Exhausted
		}

		currentValue, err := receiver.generator()
		if err != nil {
			var zeroReturnValue T
			return zeroReturnValue, err
		}

		itemsSeen++

		return currentValue, err
	}

	return newLink(limitGenerator)
}

// Distinct removes any duplicates.
func (receiver *Link[T]) Distinct() *Link[T] {
	seenTracker := helper.NewSet()

	distinctGenerator := func() (T, error) {
		//go through the generator until you find an item that hasn't been seen yet
		for {
			valueToDistinct, err := receiver.generator()
			if err != nil {
				var zeroReturnValue T
				return zeroReturnValue, err
			}

			if !seenTracker.Contains(valueToDistinct) {
				seenTracker.Add(valueToDistinct)
				return valueToDistinct, nil
			}
		}
	}

	return newLink(distinctGenerator)
}

// Flatten will iterate over all the values in the chain, but any value encountered that is a range-able container itself will also have its values iterated over first before continuing with the remaining values in the chain.  Maps flatten to its `keyvalue.KeyValuer` key and value pairs.
func (receiver *Link) Flatten() *Link {
	var currentGenerator func() (interface{}, error)

	flattenGenerator := func() (interface{}, error) {
		var innerValue interface{}
		var err error

		for innerValue == nil {
			if currentGenerator == nil {
				var currentValue interface{}
				currentValue, err = receiver.generator()
				if err != nil {
					return 0, err
				}

				if helper.IsSlice(currentValue) {
					currentGenerator = generator.FromSlice(currentValue)
				} else if helper.IsArray(currentValue) {
					currentGenerator = generator.FromArray(currentValue)
				} else if helper.IsChannel(currentValue) {
					currentGenerator = generator.FromChannel(currentValue)
				} else if helper.IsMap(currentValue) {
					currentGenerator = generator.FromMap(currentValue)
				} else {
					//it's some basic value, just return that
					innerValue = currentValue
					break
				}
			}

			innerValue, err = currentGenerator()
			if errors.Is(err, generator.Exhausted) {
				//the current generator is exhausted, set it to nil so we grab the next generator
				innerValue = nil
				currentGenerator = nil
			}
		}

		return innerValue, err
	}

	return newLink(flattenGenerator)
}

type Rangeable[T ~comparable, S any] interface {
	[]T | chan T | map[T]S
}

func FlattenFunction[T []S, S any](receiver *Link[T]) *Link[S] {
	var currentGenerator func() (S, error)

	flattenGenerator := func() (S, error) {
		var innerValue S
		var err error

		for innerValue == nil {
			if currentGenerator == nil {
				var currentValue T
				currentValue, err = receiver.generator()
				if err != nil {
					var zeroReturnValue S
					return zeroReturnValue, err
				}

				if helper.IsSlice(currentValue) {
					currentGenerator = generator.FromSlice[S](currentValue)
				} else if helper.IsArray(currentValue) {
					currentGenerator = generator.FromArray(currentValue)
				} else if helper.IsChannel(currentValue) {
					currentGenerator = generator.FromChannel(currentValue)
				} else if helper.IsMap(currentValue) {
					currentGenerator = generator.FromMap(currentValue)
				} else {
					//it's some basic value, just return that
					innerValue = currentValue
					break
				}
			}

			innerValue, err = currentGenerator()
			if errors.Is(err, generator.Exhausted) {
				//the current generator is exhausted, set it to nil so we grab the next generator
				innerValue = nil
				currentGenerator = nil
			}
		}

		return innerValue, err
	}

	return newLink(flattenGenerator)
}

// Sort sorts the chain given the `Less` function returned from the `returnLessFunction` function parameter.  The `returnLessFunction` function is called with the entire serialized chain as a slice and _returns_ a function that satisfies the same requirements as the Interface type's `Less` function (https://pkg.go.dev/sort#Interface).  This method is expensive because it must serialize all the values into a slice first.
func (receiver *Link[T]) Sort(returnLessFunction func([]T) func(int, int) bool) *Link[T] {
	serializedSlice, err := receiver.Slice()
	if err != nil {
		//there was an error during serialization, so no need to do the work of sorting
		//just always return the error that occurred
		generation := func() (T, error) {
			var zeroReturnValue T
			return zeroReturnValue, err
		}
		return newLink(generation)
	}

	lessFunction := returnLessFunction(serializedSlice)
	sort.Slice(serializedSlice, lessFunction)

	generation := generator.FromSlice(serializedSlice)

	return newLink(generation)
}

// Reverse reverses the order of the chain.  The last item will be first, and the first item will be last.  This method is expensive because it must serialize all the values into a slice first.
func (receiver *Link[T]) Reverse() *Link[T] {
	serializedSlice, err := receiver.Slice()
	if err != nil {
		//there was an error during serialization, so no need to do the work of reversing
		//just always return the error that occurred
		generation := func() (T, error) {
			var zeroReturnValue T
			return zeroReturnValue, err
		}
		return newLink(generation)
	}

	for startIndex, endIndex := 0, len(serializedSlice)-1; startIndex <= endIndex; startIndex, endIndex = startIndex+1, endIndex-1 {
		serializedSlice[startIndex], serializedSlice[endIndex] = serializedSlice[endIndex], serializedSlice[startIndex]
	}

	generation := generator.FromSlice(serializedSlice)

	return newLink(generation)
}
