package slice_chain

import (
	"github.com/halprin/slice-chain/generator"
	"github.com/halprin/slice-chain/intermediate"
)

func FromSlice(theSlice []int) *intermediate.Link {
	sliceGenerator := generator.FromSlice(theSlice)

	link := intermediate.NewLink(sliceGenerator)
	return link
}
