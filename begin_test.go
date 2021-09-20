package slice_chain

import (
	"fmt"
	"testing"
)

func TestFunStuff(t *testing.T) {
	stringSlice := []string{"DogCows", "goes", "Moof!", "Do", "you", "like", "Clarus", "the", "DogCow?"}
	chain := FromSlice(stringSlice)

	outputSlice := chain.Filter(func(value interface{}) bool {
		stringValue, ok := value.(string)
		if !ok {
			return false
		}

		return len(stringValue) % 2 == 0
	}).Skip(1).Slice()

	fmt.Println(outputSlice)
}
