package rangechain

import (
	"sort"

	"github.com/halprin/rangechain/internal/generator"
	"github.com/halprin/rangechain/internal/helper"
)

// Map runs the `mapFunction` parameter against all the values in the chain. In that function, return what you want to change the value into or an optional error.
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

// Filter runs the `filterFunction` parameter against all the values in the chain. Returning true keeps the value; returning false drops it.
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

// Skip skips the next `skipNumber` values in the chain (including any errors already in flight).
func (receiver *Link[T]) Skip(skipNumber int) *Link[T] {
	for count := 0; count < skipNumber; count++ {
		_, _ = receiver.generator()
	}

	return newLink(receiver.generator)
}

// Limit stops the chain after `keepSize` values have been emitted.
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

// DistinctFunc removes duplicates. Two values whose `keyFunction` returns the same value are considered equal. Use `func(v T) T { return v }` when the values are already comparable.
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

// Sort sorts the chain using a `Less` function returned by the `returnLessFunction` parameter. The returned function must satisfy the same requirements as the Interface type's `Less` function (https://pkg.go.dev/sort#Interface). See the TestSortingMaps example in example_test.go. Expensive because it serializes the chain first.
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

// Reverse reverses the order of the chain. Expensive because it serializes the chain first.
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
