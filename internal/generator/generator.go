// Package generator has many helper methods for generating values in a chain.
package generator

import (
	"errors"
	"github.com/halprin/rangechain/keyvalue"
	"iter"
	"maps"
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
	next, stop := iter.Pull2(maps.All(aMap))

	return func() (keyvalue.KeyValuer[K, V], error) {
		key, value, ok := next()
		if !ok {
			stop()
			return nil, Exhausted
		}

		return &mapTuple[K, V]{
			TheKey:   key,
			TheValue: value,
		}, nil
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

// FromSeq2 creates a generator for an iter.Seq2 of key/value pairs.
func FromSeq2[K, V any](seq iter.Seq2[K, V]) func() (keyvalue.KeyValuer[K, V], error) {
	next, stop := iter.Pull2(seq)

	return func() (keyvalue.KeyValuer[K, V], error) {
		key, value, ok := next()
		if !ok {
			stop()
			return nil, Exhausted
		}

		return &mapTuple[K, V]{
			TheKey:   key,
			TheValue: value,
		}, nil
	}
}
