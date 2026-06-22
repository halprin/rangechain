package rangechain

import (
	"fmt"
	"strings"
	"testing"

	"github.com/halprin/rangechain/keyvalue"
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

func TestFlattenExample(t *testing.T) {
	nested := [][]int{{1, 2, 3}, {4, 5, 6}}
	flat, _ := FromSlice(nested).Flatten[int]().Slice()
	fmt.Println(flat) // [1 2 3 4 5 6]
}

func TestDistinctExample(t *testing.T) {
	words := []string{"Moof", "moof", "DogCow", "dogcow", "Clarus"}
	distinct, _ := FromSlice(words).
		DistinctFunc(func(value string) string { return strings.ToLower(value) }).
		Slice()
	fmt.Println(distinct) // [Moof DogCow Clarus]
}

func TestIteratorExample(t *testing.T) {
	chain := FromSlice([]int{1, 2, 3, 4, 5}).
		Filter(func(value int) (bool, error) { return value%2 == 1, nil }).
		Map(func(value int) (int, error) { return value * 10, nil })

	for value, err := range chain.Iterator() {
		if err != nil {
			break
		}
		fmt.Println(value) // 10, 30, 50
	}
}

func TestFromIteratorExample(t *testing.T) {
	countdown := func(yield func(int) bool) {
		for i := 3; i > 0; i-- {
			if !yield(i) {
				return
			}
		}
	}

	launched, _ := FromIterator(countdown).
		Map(func(value int) (string, error) {
			return fmt.Sprintf("T-minus %d", value), nil
		}).
		Slice()
	fmt.Println(launched) // [T-minus 3 T-minus 2 T-minus 1]
}

func TestAllMatchExample(t *testing.T) {
	scores := []int{82, 91, 77, 100}
	allPassing, _ := FromSlice(scores).
		AllMatch(func(value int) (bool, error) { return value >= 70, nil })
	fmt.Println(allPassing) // true
}
