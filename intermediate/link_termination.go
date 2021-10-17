package intermediate

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

func (receiver *Link) Reduce(reduceFunction func(interface{}, interface{}) interface{}) *interface{} {
	nextItem, err := receiver.generator()
	if err != nil {
		return nil
	}

	intermediateItem, err := receiver.generator()
	if err != nil {
		return &nextItem
	}

	for err == nil {
		intermediateItem = reduceFunction(intermediateItem, nextItem)
		nextItem, err = receiver.generator()
	}

	return &intermediateItem
}

func (receiver *Link) ReduceWithInitialValue(reduceFunction func(interface{}, interface{}) interface{}, initialValue interface{}) interface{} {
	nextItem, err := receiver.generator()
	if err != nil {
		return initialValue
	}

	intermediateItem := initialValue

	for err == nil {
		intermediateItem = reduceFunction(intermediateItem, nextItem)
		nextItem, err = receiver.generator()
	}

	return intermediateItem
}
