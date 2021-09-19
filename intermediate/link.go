package intermediate

type Link struct {
	generator func() int
}

func NewLink(generator func() int) *Link {
	return &Link{
		generator: generator,
	}
}

//chain methods

func (receiver *Link) filter(filterFunction func(int) bool) *Link {
	filterGenerator := func() int {
		//go through the generator until you find an item that stays
		for {
			valueToFilter := receiver.generator()
			if filterFunction(valueToFilter) {
				return valueToFilter
			}
		}
	}

	return NewLink(filterGenerator)
}

//termination methods

func (receiver *Link) slice() (endSlice []int) {
	//var endSlice []int

	defer func() {
		recover()
		//if p != nil {
		//	return endSlice
		//}
	}()

	for {
		currentValue := receiver.generator()
		endSlice = append(endSlice, currentValue)
	}

	return endSlice
}
