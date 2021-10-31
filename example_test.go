package rangechain

import (
	"fmt"
	"github.com/halprin/rangechain/internal/generator"
	"testing"
)

func TestFunStuff(t *testing.T) {
	stringSlice := []string{"DogCows", "goes", "Moof", "Do", "you", "like", "Clarus", "the", "DogCow"}
	chain := FromSlice(stringSlice)

	outputSlice, _ := chain.
		Map(func(value interface{}) (interface{}, error) {
			stringValue := value.(string)
			return stringValue + " not", nil
		}).
		Filter(func(value interface{}) (bool, error) {
			stringValue := value.(string)

			return len(stringValue) % 2 == 0, nil
		}).

		Skip(1).
		Slice()

	fmt.Println(outputSlice)
}

func TestReduceToMapExample(t *testing.T) {
	stringSlice := []string{"DogCows", "goes", "Moof", "Do", "you", "like", "Clarus", "the", "DogCow"}
	chain := FromSlice(stringSlice)

	outputMap, _ := chain.
		ReduceWithInitialValue(func(firstItem interface{}, secondItem interface{}) (interface{}, error) {
			reductionMap := firstItem.(map[string]int)
			stringItem := secondItem.(string)

			reductionMap[stringItem] = len(stringItem)

			return reductionMap, nil
		}, map[string]int{})

	fmt.Println(outputMap)
}

func TestSortingMaps(t *testing.T) {
	aMap := map[string]int{
		"DogCow": 10,
		"System 7": 7,
		"Mac OS 8": 8,
		"Mac OS 9": 9,
		"Mac OS X": 10,
		"QuickTime": 3,
		"Exposé": 7,
		"Control Strip": 6,
		"Finder": 5,
	}

	chain := FromMap(aMap)
	sortedAppleStuff, _ := chain.Sort(func(mapValuesToSort []interface{}) func(int, int) bool {
		return func(index1 int, index2 int) bool {
			mapValue1 := mapValuesToSort[index1].(generator.MapTuple)
			mapValue2 := mapValuesToSort[index2].(generator.MapTuple)

			rating1 := mapValue1.Value.(int)
			rating2 := mapValue2.Value.(int)

			return rating1 > rating2
		}
	}).Map(func(value interface{}) (interface{}, error) {
		mapValue := value.(generator.MapTuple)
		return mapValue.Key, nil
	}).Slice()

	fmt.Println(sortedAppleStuff)
}
