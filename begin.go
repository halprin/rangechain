package slice_chain

import (
	channelGeneratorImport "github.com/halprin/slice-chain/channelGenerator"
	"github.com/halprin/slice-chain/channelIntermediate"
	"github.com/halprin/slice-chain/generator"
	"github.com/halprin/slice-chain/helper"
	"github.com/halprin/slice-chain/intermediate"
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

func FromMap(aMap map[interface{}]interface{}) *intermediate.Link {
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
