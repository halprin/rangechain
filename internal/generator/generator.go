// Package generator has many helper methods for generating values in a chain.
package generator

import (
	"errors"
	"github.com/halprin/rangechain/keyvalue"
	"iter"
)

// Exhausted is returned as an expected error from the generators to designate an end of the generator.
var Exhausted = errors.New("generator exhausted")

// FromSlice creates a generator for a slice.
func FromSlice[T any](slice []T) func() (T, error) {
	currentIndex := 0

	return func() (T, error) {
		if currentIndex >= len(slice) {
			var zero T
			return zero, Exhausted
		}

		value := slice[currentIndex]
		currentIndex++

		return value, nil
	}
}

// FromChannel creates a generator for a channel.
func FromChannel[T any](channel <-chan T) func() (T, error) {
	return func() (T, error) {
		value, ok := <-channel
		if !ok {
			var zero T
			return zero, Exhausted
		}

		return value, nil
	}
}

// FromMap creates a generator for a map.
func FromMap[K comparable, V any](aMap map[K]V) func() (keyvalue.KeyValuer[K, V], error) {
	keys := make([]K, 0, len(aMap))
	for key := range aMap {
		keys = append(keys, key)
	}

	currentIndex := 0

	return func() (keyvalue.KeyValuer[K, V], error) {
		if currentIndex >= len(keys) {
			return nil, Exhausted
		}

		key := keys[currentIndex]
		currentIndex++

		tuple := &mapTuple[K, V]{
			TheKey:   key,
			TheValue: aMap[key],
		}

		return tuple, nil
	}
}

// FromIterator creates a generator for an iter.Seq.
func FromIterator[T any](anIterator iter.Seq[T]) func() (T, error) {
	next, stop := iter.Pull(anIterator)

	return func() (T, error) {
		value, ok := next()
		if !ok {
			stop()
			var zero T
			return zero, Exhausted
		}

		return value, nil
	}
}
