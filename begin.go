package slice_chain

import (
	"github.com/halprin/slice-chain/channelGenerator"
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

func FromSliceWithChannels(slice interface{}) *channelIntermediate.Link {
	interfaceSlice := helper.InterfaceSlice(slice)
	sliceGenerator := channelGenerator.FromSlice(interfaceSlice)

	link := channelIntermediate.NewLink(sliceGenerator)
	return link
}
