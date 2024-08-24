package rangechain

import (
	"github.com/halprin/rangechain/internal/generator"
	"iter"
)

// FromSlice starts the chain with the supplied slice.
// Chaining and terminating methods can now be called on the result.
func FromSlice(slice interface{}) *Link {
	sliceGenerator := generator.FromSlice(slice)

	link := newLink(sliceGenerator)
	return link
}

// FromArray starts the chain with the supplied array.
// Chaining and terminating methods can now be called on the result.
func FromArray(array interface{}) *Link {
	arrayGenerator := generator.FromArray(array)

	link := newLink(arrayGenerator)
	return link
}

// FromChannel starts the chain with the supplied channel.
// Chaining and terminating methods can now be called on the result.
func FromChannel(channel interface{}) *Link {
	channelGenerator := generator.FromChannel(channel)

	link := newLink(channelGenerator)
	return link
}

// FromMap starts the chain with the supplied map.
// Chaining and terminating methods can now be called on the result.  The singular value used to represent the key and value pairs is `keyvalue.KeyValuer` of `github.com/halprin/rangechain/keyvalue`.
func FromMap(aMap interface{}) *Link {
	mapGenerator := generator.FromMap(aMap)

	link := newLink(mapGenerator)
	return link
}

// FromIterator starts the chain with the supplied iterator.
// Chaining and terminating methods can now be called on the result.
func FromIterator[T any](anIterator iter.Seq[T]) *Link {
	iterGenerator := generator.FromIterator(anIterator)

	link := newLink(iterGenerator)
	return link
}
