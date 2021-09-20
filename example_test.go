package slice_chain

import (
	"fmt"
	"testing"
)

func TestFunStuff(t *testing.T) {
	stringSlice := []string{"DogCows", "goes", "Moof", "Do", "you", "like", "Clarus", "the", "DogCow"}
	chain := FromSlice(stringSlice)

	outputSlice := chain.
		Map(func(value interface{}) interface{} {
			stringValue := value.(string)
			return stringValue + " not"
		}).
		Filter(func(value interface{}) bool {
			stringValue := value.(string)

			return len(stringValue) % 2 == 0
		}).

		Skip(1).
		Slice()

	fmt.Println(outputSlice)
}