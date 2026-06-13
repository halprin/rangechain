package rangechain

import (
	"github.com/halprin/rangechain/internal/generator"
	"github.com/halprin/rangechain/keyvalue"
	"iter"
)

// FromSlice starts the chain with the supplied slice.
// Chaining and terminating methods can now be called on the result.
func FromSlice[T any](slice []T) *Link[T] {
	return newLink(generator.FromSlice(slice))
}

// FromChannel starts the chain with the supplied channel.
// Chaining and terminating methods can now be called on the result.
func FromChannel[T any](channel <-chan T) *Link[T] {
	return newLink(generator.FromChannel(channel))
}

// FromMap starts the chain with the supplied map.
// Chaining and terminating methods can now be called on the result.  The singular value used to represent the key and value pairs is `keyvalue.KeyValuer` of `github.com/halprin/rangechain/keyvalue`.
func FromMap[K comparable, V any](aMap map[K]V) *Link[keyvalue.KeyValuer[K, V]] {
	return newLink(generator.FromMap(aMap))
}

// FromIterator starts the chain with the supplied iterator.
// Chaining and terminating methods can now be called on the result.
func FromIterator[T any](anIterator iter.Seq[T]) *Link[T] {
	return newLink(generator.FromIterator(anIterator))
}
