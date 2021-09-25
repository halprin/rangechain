package channelIntermediate

import (
	"github.com/halprin/slice-chain/helper"
)

type Link struct {
	generator chan interface{}
}

func NewLink(generator chan interface{}) *Link {
	return &Link{
		generator: generator,
	}
}

//chain methods

func (receiver *Link) Map(mapFunction func(interface{}) interface{}) *Link {
	mapGenerator := make(chan interface{})

	go func() {
		for currentValue := range receiver.generator {
			mappedValue := mapFunction(currentValue)
			mapGenerator <- mappedValue
		}

		close(mapGenerator)
	}()


	return NewLink(mapGenerator)
}

func (receiver *Link) Filter(filterFunction func(interface{}) bool) *Link {
	filterGenerator := make(chan interface{})

	go func() {
		for currentValue := range receiver.generator {
			keepValue := filterFunction(currentValue)
			if keepValue {
				filterGenerator <- currentValue
			}
		}

		close(filterGenerator)
	}()

	return NewLink(filterGenerator)
}

func (receiver *Link) Skip(skipNumber int) *Link {
	for count := 0; count < skipNumber; count++ {
		<-receiver.generator
	}

	return NewLink(receiver.generator)
}

func (receiver *Link) Limit(keepSize int) *Link {
	limitGenerator := make(chan interface{})

	go func() {
		itemsSeen := 0
		for currentValue := range receiver.generator {
			if itemsSeen >= keepSize {
				break
			}

			limitGenerator <- currentValue

			itemsSeen++
		}

		close(limitGenerator)
	}()

	return NewLink(limitGenerator)
}

func (receiver *Link) Distinct() *Link {
	distinctGenerator := make(chan interface{})

	go func() {
		seenTracker := helper.NewSet()
		for currentValue := range receiver.generator {
			if !seenTracker.Contains(currentValue) {
				seenTracker.Add(currentValue)
				distinctGenerator <- currentValue
			}
		}

		close(distinctGenerator)
	}()

	return NewLink(distinctGenerator)
}

//func (receiver *Link) Flatten() *Link {
//	var currentGenerator func() (interface{}, error)
//
//	flattenGenerator := func() (interface{}, error) {
//		var innerValue interface{}
//		var err error
//
//		for innerValue == nil {
//			if currentGenerator == nil {
//				var currentValue interface{}
//				currentValue, err = receiver.generator()
//				if err != nil {
//					return 0, err
//				}
//
//				if helper.IsSlice(currentValue) {
//					sliceCurrentValue := helper.InterfaceSlice(currentValue)
//					currentGenerator = generator.FromSlice(sliceCurrentValue)
//				} else {
//					//it's some basic value, just return that
//					innerValue = currentValue
//					break
//				}
//			}
//
//			innerValue, err = currentGenerator()
//			if err != nil {
//				//the current generator is exhausted, set it to nil so we grab the next generator
//				innerValue = nil
//				currentGenerator = nil
//			}
//		}
//
//		return innerValue, err
//	}
//
//	return NewLink(flattenGenerator)
//}
//
//func (receiver *Link) Sort(returnLessFunction func([]interface{}) func(int, int) bool) *Link {
//	serializedSlice := receiver.Slice()
//
//	lessFunction := returnLessFunction(serializedSlice)
//	sort.Slice(serializedSlice, lessFunction)
//
//
//	generation := generator.FromSlice(serializedSlice)
//
//	return NewLink(generation)
//}
//
//func (receiver *Link) Reverse() *Link {
//	serializedSlice := receiver.Slice()
//
//	for startIndex, endIndex := 0, len(serializedSlice) - 1; startIndex <= endIndex; startIndex, endIndex = startIndex + 1, endIndex - 1 {
//		serializedSlice[startIndex], serializedSlice[endIndex] = serializedSlice[endIndex], serializedSlice[startIndex]
//	}
//
//	generation := generator.FromSlice(serializedSlice)
//
//	return NewLink(generation)
//}

//termination methods

func (receiver *Link) Slice() []interface{} {
	endSlice := []interface{}{}

	for currentValue := range receiver.generator {
		endSlice = append(endSlice, currentValue)
	}

	return endSlice
}

func (receiver *Link) ForEach(forEachFunction func(interface{})) {
	for currentValue := range receiver.generator {
		forEachFunction(currentValue)
	}
}

func (receiver *Link) Count() int {
	count := 0

	for range receiver.generator {
		count++
	}

	return count
}

//func (receiver *Link) First() *interface{} {
//	value, err := receiver.generator()
//	if err != nil {
//		return nil
//	}
//
//	return &value
//}

func (receiver *Link) Last() *interface{} {
	var last *interface{}

	for currentValue := range receiver.generator {
		last = &currentValue
	}

	return last
}

func (receiver *Link) AllMatch(allMatchFunction func(interface{}) bool) bool {
	for currentValue := range receiver.generator {
		if !allMatchFunction(currentValue) {
			return false
		}
	}

	return true
}

func (receiver *Link) AnyMatch(anyMatchFunction func(interface{}) bool) bool {
	for currentValue := range receiver.generator {
		if anyMatchFunction(currentValue) {
			return true
		}
	}

	return false
}

func (receiver *Link) NoneMatch(noneMatchFunction func(interface{}) bool) bool {
	return !receiver.AnyMatch(noneMatchFunction)
}

//func (receiver *Link) Reduce(reduceFunction func(interface{}, interface{}) interface{}) *interface{} {
//	nextItem, err := receiver.generator()
//	if err != nil {
//		return nil
//	}
//
//	intermediateItem, err := receiver.generator()
//	if err != nil {
//		return &nextItem
//	}
//
//	for err == nil {
//		intermediateItem = reduceFunction(intermediateItem, nextItem)
//		nextItem, err = receiver.generator()
//	}
//
//	return &intermediateItem
//}
//
//func (receiver *Link) ReduceWithInitialValue(reduceFunction func(interface{}, interface{}) interface{}, initialValue interface{}) interface{} {
//	nextItem, err := receiver.generator()
//	if err != nil {
//		return initialValue
//	}
//
//	intermediateItem := initialValue
//
//	for err == nil {
//		intermediateItem = reduceFunction(intermediateItem, nextItem)
//		nextItem, err = receiver.generator()
//	}
//
//	return intermediateItem
//}
