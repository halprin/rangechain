package intermediate

import (
	"github.com/halprin/slice-chain/generator"
	"github.com/halprin/slice-chain/helper"
	"sort"
)

type Link struct {
	generator func() (interface{}, error)
}

func NewLink(generator func() (interface{}, error)) *Link {
	return &Link{
		generator: generator,
	}
}

//chain methods

func (receiver *Link) Map(mapFunction func(interface{}) interface{}) *Link {
	mapGenerator := func() (interface{}, error) {
		valueToMap, err := receiver.generator()
		if err != nil {
			return 0, err
		}

		return mapFunction(valueToMap), nil
	}

	return NewLink(mapGenerator)
}

func (receiver *Link) Filter(filterFunction func(interface{}) bool) *Link {
	filterGenerator := func() (interface{}, error) {
		//go through the generator until you find an item that stays
		for {
			valueToFilter, err := receiver.generator()
			if err != nil {
				return 0, err
			}

			if filterFunction(valueToFilter) {
				return valueToFilter, nil
			}
		}
	}

	return NewLink(filterGenerator)
}

func (receiver *Link) Skip(skipNumber int) *Link {
	for count := 0; count < skipNumber; count++ {
		_, _ = receiver.generator()
	}

	return NewLink(receiver.generator)
}

func (receiver *Link) Limit(keepSize int) *Link {
	itemsSeen := 0

	limitGenerator := func() (interface{}, error) {
		if itemsSeen >= keepSize {
			return 0, generator.Exhausted
		}

		currentValue, err := receiver.generator()
		if err != nil {
			return 0, err
		}

		itemsSeen++

		return currentValue, err
	}

	return NewLink(limitGenerator)
}

func (receiver *Link) Distinct() *Link {
	seenTracker := helper.NewSet()

	distinctGenerator := func() (interface{}, error) {
		//go through the generator until you find an item that hasn't been seen yet
		for {
			valueToDistinct, err := receiver.generator()
			if err != nil {
				return 0, err
			}

			if !seenTracker.Contains(valueToDistinct) {
				seenTracker.Add(valueToDistinct)
				return valueToDistinct, nil
			}
		}
	}

	return NewLink(distinctGenerator)
}

func (receiver *Link) Flatten() *Link {
	var currentGenerator func() (interface{}, error)

	flattenGenerator := func() (interface{}, error) {
		var innerValue interface{}
		var err error

		for innerValue == nil {
			if currentGenerator == nil {
				var currentValue interface{}
				currentValue, err = receiver.generator()
				if err != nil {
					return 0, err
				}

				if helper.IsSlice(currentValue) {
					sliceCurrentValue := helper.InterfaceSlice(currentValue)
					currentGenerator = generator.FromSlice(sliceCurrentValue)
				} else {
					//it's some basic value, just return that
					innerValue = currentValue
					break
				}
			}

			innerValue, err = currentGenerator()
			if err != nil {
				//the current generator is exhausted, set it to nil so we grab the next generator
				innerValue = nil
				currentGenerator = nil
			}
		}

		return innerValue, err
	}

	return NewLink(flattenGenerator)
}

func (receiver *Link) Sort(returnLessFunction func([]interface{}) func(int, int) bool) *Link {
	serializedSlice := receiver.Slice()

	lessFunction := returnLessFunction(serializedSlice)
	sort.Slice(serializedSlice, lessFunction)


	generation := generator.FromSlice(serializedSlice)

	return NewLink(generation)
}

//termination methods

func (receiver *Link) Slice() []interface{} {
	endSlice := []interface{}{}

	for {
		currentValue, err := receiver.generator()
		if err != nil {
			return endSlice
		}

		endSlice = append(endSlice, currentValue)
	}
}

func (receiver *Link) ForEach(forEachFunction func(interface{})) {
	for {
		currentValue, err := receiver.generator()
		if err != nil {
			return
		}

		forEachFunction(currentValue)
	}
}

func (receiver *Link) Count() int {
	count := 0
	for {
		_, err := receiver.generator()
		if err != nil {
			return count
		}

		count++
	}
}

func (receiver *Link) First() *interface{} {
	value, err := receiver.generator()
	if err != nil {
		return nil
	}

	return &value
}

func (receiver *Link) Last() *interface{} {
	var last *interface{}

	for {
		currentValue, err := receiver.generator()
		if err != nil {
			return last
		}

		last = &currentValue
	}
}

func (receiver *Link) AllMatch(allMatchFunction func(interface{}) bool) bool {
	for {
		currentValue, err := receiver.generator()
		if err != nil {
			return true
		}

		if !allMatchFunction(currentValue) {
			return false
		}
	}
}

func (receiver *Link) AnyMatch(anyMatchFunction func(interface{}) bool) bool {
	for {
		currentValue, err := receiver.generator()
		if err != nil {
			return false
		}

		if anyMatchFunction(currentValue) {
			return true
		}
	}
}

func (receiver *Link) NoneMatch(noneMatchFunction func(interface{}) bool) bool {
	return !receiver.AnyMatch(noneMatchFunction)
}
