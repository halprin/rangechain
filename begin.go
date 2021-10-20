package slice_chain

import (
	channelGeneratorImport "github.com/halprin/rangechain/channelGenerator"
	"github.com/halprin/rangechain/channelIntermediate"
	"github.com/halprin/rangechain/generator"
	"github.com/halprin/rangechain/helper"
	"github.com/halprin/rangechain/intermediate"
)

func FromSlice(slice interface{}) *intermediate.Link {
	sliceGenerator := generator.FromSlice(slice)

	link := intermediate.NewLink(sliceGenerator)
	return link
}

func FromArray(array interface{}) *intermediate.Link {
	arrayGenerator := generator.FromArray(array)

	link := intermediate.NewLink(arrayGenerator)
	return link
}

func FromChannel(channel interface{}) *intermediate.Link {
	channelGenerator := generator.FromChannel(channel)

	link := intermediate.NewLink(channelGenerator)
	return link
}

func FromMap(aMap interface{}) *intermediate.Link {
	mapGenerator := generator.FromMap(aMap)

	link := intermediate.NewLink(mapGenerator)
	return link
}

func FromSliceWithChannels(slice interface{}) *channelIntermediate.Link {
	interfaceSlice := helper.InterfaceSlice(slice)
	sliceGenerator := channelGeneratorImport.FromSlice(interfaceSlice)

	link := channelIntermediate.NewLink(sliceGenerator)
	return link
}
