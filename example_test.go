package rangechain

import (
	"fmt"
	"github.com/halprin/rangechain/keyvalue"
	"testing"
)

func TestFunStuff(t *testing.T) {
	stringSlice := []string{"DogCows", "goes", "Moof", "Do", "you", "like", "Clarus", "the", "DogCow"}
	chain := FromSlice(stringSlice)

	outputSlice, _ := chain.
		Map(func(value string) (string, error) {
			return value + " not", nil
		}).
		Filter(func(value string) (bool, error) {
			return len(value)%2 == 0, nil
		}).
		Skip(1).
		Slice()

	fmt.Println(outputSlice)
}

func TestReduceToMapExample(t *testing.T) {
	stringSlice := []string{"DogCows", "goes", "Moof", "Do", "you", "like", "Clarus", "the", "DogCow"}
	chain := FromSlice(stringSlice)

	outputMap, _ := chain.
		ReduceWithInitialValue(func(reductionMap map[string]int, item string) (map[string]int, error) {
			reductionMap[item] = len(item)
			return reductionMap, nil
		}, map[string]int{})

	fmt.Println(outputMap)
}

func TestSortingMaps(t *testing.T) {
	aMap := map[string]int{
		"DogCow":        10,
		"System 7":      7,
		"Mac OS 8":      8,
		"Mac OS 9":      9,
		"Mac OS X":      10,
		"QuickTime":     3,
		"Exposé":        7,
		"Control Strip": 6,
		"Finder":        5,
	}

	chain := FromMap(aMap)
	sortedAppleStuff, _ := chain.Sort(func(mapValuesToSort []keyvalue.KeyValuer[string, int]) func(int, int) bool {
		return func(index1 int, index2 int) bool {
			return mapValuesToSort[index1].Value() > mapValuesToSort[index2].Value()
		}
	}).Map(func(value keyvalue.KeyValuer[string, int]) (string, error) {
		return value.Key(), nil
	}).Slice()

	fmt.Println(sortedAppleStuff)
}
