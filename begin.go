package slice_chain

import (
	"github.com/halprin/slice-chain/generator"
	"github.com/halprin/slice-chain/helper"
	"github.com/halprin/slice-chain/intermediate"
)

func FromSlice(slice interface{}) *intermediate.Link {
	interfaceSlice := helper.InterfaceSlice(slice)
	sliceGenerator := generator.FromSlice(interfaceSlice)

	link := intermediate.NewLink(sliceGenerator)
	return link
}
