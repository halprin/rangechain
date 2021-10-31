package rangechain

import (
	"github.com/halprin/rangechain/generator"
	"github.com/halprin/rangechain/intermediate"
)

// FromSlice starts the chain with the supplied slice.
// Chaining and terminating methods can now be called on the result.
func FromSlice(slice interface{}) *intermediate.Link {
	sliceGenerator := generator.FromSlice(slice)

	link := intermediate.NewLink(sliceGenerator)
	return link
}

// FromArray starts the chain with the supplied array.
// Chaining and terminating methods can now be called on the result.
func FromArray(array interface{}) *intermediate.Link {
	arrayGenerator := generator.FromArray(array)

	link := intermediate.NewLink(arrayGenerator)
	return link
}

// FromChannel starts the chain with the supplied channel.
// Chaining and terminating methods can now be called on the result.
func FromChannel(channel interface{}) *intermediate.Link {
	channelGenerator := generator.FromChannel(channel)

	link := intermediate.NewLink(channelGenerator)
	return link
}

// FromMap starts the chain with the supplied map.
// Chaining and terminating methods can now be called on the result.  The singular value used to represent the key and value pairs is `generator.MapTuple` of `github.com/halprin/rangechain/generator`.
func FromMap(aMap interface{}) *intermediate.Link {
	mapGenerator := generator.FromMap(aMap)

	link := intermediate.NewLink(mapGenerator)
	return link
}
