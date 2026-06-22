package rangechain

import (
	"sort"

	"github.com/halprin/rangechain/internal/generator"
	"github.com/halprin/rangechain/internal/helper"
)

// Map will run the `mapFunction` parameter function against all the values in the chain.  In that function, return what you want to change the value into or an optional error if an error is encountered.
func (receiver *Link[T]) Map[U any](mapFunction func(T) (U, error)) *Link[U] {
	mapGenerator := func() (U, error) {
		valueToMap, err := receiver.generator()
		if err != nil {
			var zero U
			return zero, err
		}

		return mapFunction(valueToMap)
	}

	return newLink(mapGenerator)
}

// Filter will run the `filterFunction` parameter function against all the values in the chain.  In that function, on return of true, the value will stay, or on false, the value will be removed.
func (receiver *Link[T]) Filter(filterFunction func(T) (bool, error)) *Link[T] {
	filterGenerator := func() (T, error) {
		for {
			valueToFilter, err := receiver.generator()
			if err != nil {
				var zero T
				return zero, err
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
			var zero T
			return zero, generator.Exhausted
		}

		currentValue, err := receiver.generator()
		if err != nil {
			var zero T
			return zero, err
		}

		itemsSeen++

		return currentValue, err
	}

	return newLink(limitGenerator)
}

// DistinctFunc removes any duplicates as identified by the parameter `keyFunction`.  Two values whose `keyFunction` returns the same `K` are considered equal.
func (receiver *Link[T]) DistinctFunc[K comparable](keyFunction func(T) K) *Link[T] {
	seenTracker := helper.NewSet[K]()

	distinctGenerator := func() (T, error) {
		for {
			valueToDistinct, err := receiver.generator()
			if err != nil {
				var zero T
				return zero, err
			}

			key := keyFunction(valueToDistinct)
			if !seenTracker.Contains(key) {
				seenTracker.Add(key)
				return valueToDistinct, nil
			}
		}
	}

	return newLink(distinctGenerator)
}

// Sort sorts the chain given the `Less` function returned from the `returnLessFunction` function parameter.  The `returnLessFunction` function is called with the entire serialized chain as a slice and _returns_ a function that satisfies the same requirements as the Interface type's `Less` function (https://pkg.go.dev/sort#Interface).  This method is expensive because it must serialize all the values into a slice first.
func (receiver *Link[T]) Sort(returnLessFunction func([]T) func(int, int) bool) *Link[T] {
	serializedSlice, err := receiver.Slice()
	if err != nil {
		generation := func() (T, error) {
			var zero T
			return zero, err
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
		generation := func() (T, error) {
			var zero T
			return zero, err
		}
		return newLink(generation)
	}

	for startIndex, endIndex := 0, len(serializedSlice)-1; startIndex <= endIndex; startIndex, endIndex = startIndex+1, endIndex-1 {
		serializedSlice[startIndex], serializedSlice[endIndex] = serializedSlice[endIndex], serializedSlice[startIndex]
	}

	generation := generator.FromSlice(serializedSlice)

	return newLink(generation)
}
