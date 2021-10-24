package intermediate

import (
	"errors"
	"github.com/halprin/rangechain/generator"
)

func (receiver *Link) Slice() ([]interface{}, error) {
	endSlice := []interface{}{}

	for {
		currentValue, err := receiver.generator()
		if err != nil {
			if errors.Is(err, generator.Exhausted) {
				return endSlice, nil
			} else {
				return nil, err
			}
		}

		endSlice = append(endSlice, currentValue)
	}
}

func (receiver *Link) Channel() (<-chan interface{}, <-chan error) {
	endChannel := make(chan interface{})
	errorChannel := make(chan error)

	go func() {
		for {
			currentValue, err := receiver.generator()
			if err != nil {
				//close the end channel no matter what
				//if this is due to a user error (and not Exhausted), I also want to close the endChannel first before writing to the error channel
				close(endChannel)

				if !errors.Is(err, generator.Exhausted) {
					errorChannel <- err
				}
				close(errorChannel)

				return
			}

			endChannel <- currentValue
		}
	}()

	return endChannel, errorChannel
}

func (receiver *Link) ForEach(forEachFunction func(interface{})) error {
	for {
		currentValue, err := receiver.generator()
		if err != nil {
			if errors.Is(err, generator.Exhausted) {
				return nil
			} else if !errors.Is(err, generator.Exhausted) {
				return err
			}
		}

		forEachFunction(currentValue)
	}
}

func (receiver *Link) ForEachParallel(forEachFunction func(interface{})) error {
	for {
		currentValue, err := receiver.generator()
		if err != nil {
			if errors.Is(err, generator.Exhausted) {
				return nil
			} else if !errors.Is(err, generator.Exhausted) {
				return err
			}
		}

		go forEachFunction(currentValue)
	}
}

func (receiver *Link) Count() (int, error) {
	count := 0
	var firstError error
	for {
		_, err := receiver.generator()
		if err != nil {
			if errors.Is(err, generator.Exhausted) {
				return count, firstError
			} else if !errors.Is(err, generator.Exhausted) {
				firstError = err
			}
		}

		count++
	}
}

func (receiver *Link) First() (*interface{}, error) {
	value, err := receiver.generator()
	if err != nil {
		if errors.Is(err, generator.Exhausted) {
			return nil, nil
		} else if !errors.Is(err, generator.Exhausted) {
			return nil, err
		}
	}

	return &value, nil
}

func (receiver *Link) Last() (*interface{}, error) {
	var lastValue *interface{}
	var lastError error

	for {
		currentValue, err := receiver.generator()
		if err != nil && errors.Is(err, generator.Exhausted) {
			return lastValue, lastError
		}

		lastValue = &currentValue
		lastError = err
	}
}

func (receiver *Link) AllMatch(allMatchFunction func(interface{}) (bool, error)) (bool, error) {
	for {
		currentValue, err := receiver.generator()
		if err != nil {
			if errors.Is(err, generator.Exhausted) {
				return true, nil
			} else if !errors.Is(err, generator.Exhausted) {
				return false, err
			}
		}

		match, err := allMatchFunction(currentValue)
		if err != nil {
			return false, err
		} else if !match {
			return false, nil
		}
	}
}

func (receiver *Link) AnyMatch(anyMatchFunction func(interface{}) bool) (bool, error) {
	for {
		currentValue, err := receiver.generator()
		if err != nil {
			if errors.Is(err, generator.Exhausted) {
				//we've reached the end and apparently never returned until now, so nothing matched
				return false, nil
			} else if !errors.Is(err, generator.Exhausted) {
				//we've reached an error, and never returned earlier, so nothing matched
				return false, err
			}
		}

		if anyMatchFunction(currentValue) {
			return true, nil
		}
	}
}

func (receiver *Link) NoneMatch(noneMatchFunction func(interface{}) bool) bool {
	match, _ := receiver.AnyMatch(noneMatchFunction)
	return !match
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
