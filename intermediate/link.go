package intermediate

import "github.com/halprin/slice-chain/generator"

type Link struct {
	generator func() (int, error)
}

func NewLink(generator func() (int, error)) *Link {
	return &Link{
		generator: generator,
	}
}

//chain methods

func (receiver *Link) Filter(filterFunction func(int) bool) *Link {
	filterGenerator := func() (int, error) {
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

	limitGenerator := func() (int, error) {
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

//termination methods

func (receiver *Link) Slice() []int {
	endSlice := []int{}

	for {
		currentValue, err := receiver.generator()
		if err != nil {
			return endSlice
		}

		endSlice = append(endSlice, currentValue)
	}
}

func (receiver *Link) ForEach(forEachFunction func(int)) {
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

func (receiver *Link) First() *int {
	value, err := receiver.generator()
	if err != nil {
		return nil
	}

	return &value
}

func (receiver *Link) Last() *int {
	var last *int

	for {
		currentValue, err := receiver.generator()
		if err != nil {
			return last
		}

		last = &currentValue
	}
}

func (receiver *Link) AllMatch(allMatchFunction func(int) bool) bool {
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

func (receiver *Link) AnyMatch(anyMatchFunction func(int) bool) bool {
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

func (receiver *Link) NoneMatch(noneMatchFunction func(int) bool) bool {
	return !receiver.AnyMatch(noneMatchFunction)
}
