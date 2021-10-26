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
				return endSlice, err
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
				if firstError == nil {
					firstError = err
				}
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

func (receiver *Link) AnyMatch(anyMatchFunction func(interface{}) (bool, error)) (bool, error) {
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

		match, err := anyMatchFunction(currentValue)
		if err != nil {
			return false, err
		} else if match {
			return true, nil
		}
	}
}

func (receiver *Link) NoneMatch(noneMatchFunction func(interface{}) (bool, error)) (bool, error) {
	match, err := receiver.AnyMatch(noneMatchFunction)
	return !match, err
}

func (receiver *Link) Reduce(reduceFunction func(interface{}, interface{}) (interface{}, error)) (*interface{}, error) {
	nextItem, err := receiver.generator()
	if err != nil {
		if errors.Is(err, generator.Exhausted) {
			return nil, nil
		} else if !errors.Is(err, generator.Exhausted) {
			return nil, err
		}
	}

	intermediateItem, err := receiver.generator()
	if err != nil {
		if errors.Is(err, generator.Exhausted) {
			return &nextItem, nil
		} else if !errors.Is(err, generator.Exhausted) {
			return &nextItem, err
		}
	}

	for err == nil {
		intermediateItem, err = reduceFunction(intermediateItem, nextItem)
		if err != nil {
			break
		}

		nextItem, err = receiver.generator()
	}

	if errors.Is(err, generator.Exhausted) {
		//if the error that stopped the for loop, don't report it as an error
		err = nil
	}

	return &intermediateItem, err
}

func (receiver *Link) ReduceWithInitialValue(reduceFunction func(interface{}, interface{}) (interface{}, error), initialValue interface{}) (interface{}, error) {
	nextItem, err := receiver.generator()
	if err != nil {
		if errors.Is(err, generator.Exhausted) {
			return initialValue, nil
		} else if !errors.Is(err, generator.Exhausted) {
			return initialValue, err
		}
	}

	intermediateItem := initialValue

	for err == nil {
		intermediateItem, err = reduceFunction(intermediateItem, nextItem)
		if err != nil {
			break
		}

		nextItem, err = receiver.generator()
	}

	if errors.Is(err, generator.Exhausted) {
		//if the error that stopped the for loop, don't report it as an error
		err = nil
	}

	return intermediateItem, err
}
