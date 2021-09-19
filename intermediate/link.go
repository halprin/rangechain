package intermediate

type Link struct {
	generator func() (int, error)
}

func NewLink(generator func() (int, error)) *Link {
	return &Link{
		generator: generator,
	}
}

//chain methods

func (receiver *Link) filter(filterFunction func(int) bool) *Link {
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

//termination methods

func (receiver *Link) slice() []int {
	var endSlice []int

	for {
		currentValue, err := receiver.generator()
		if err != nil {
			return endSlice
		}

		endSlice = append(endSlice, currentValue)
	}
}
