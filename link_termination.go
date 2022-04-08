package rangechain

import (
	"errors"
	"github.com/halprin/rangechain/internal/generator"
)

// Slice serializes the chain into a slice and returns it.  Also returns an error if any previous chain method generated an error.  If an error is returned, the slice is filled in until the error was encountered.
func (receiver *Link[T]) Slice() ([]T, error) {
	endSlice := []T{}

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

// Channel serializes the chain into a channel.  Also returns any errors in a channel if any previous chain method generated an error.  If an error is returned, the value channel is closed, the error is sent on the error channel, and the error channel is closed.
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

// ForEach will run the `forEachFunction` parameter function across all the values in the chain.  Also returns an error if any previous chain method generated an error.  If an error is encountered, the function stops executing against the remaining chain.
func (receiver *Link[T]) ForEach(forEachFunction func(T)) error {
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

// ForEachParallel will run the `forEachFunction` parameter function across all the values in the chain in parallel.  Also returns an error if any previous chain method generated an error.  If an error is encountered, the function stops executing against the remaining chain.  There is overhead to running in parallel so benchmark to ensure you benefit from this version.
func (receiver *Link[T]) ForEachParallel(forEachFunction func(T)) error {
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

// Count returns the number of values in the chain.  Also returns an error if any previous chain method generated an error.  Count returns an accurate number even if an error is encountered.
func (receiver *Link[T]) Count() (int, error) {
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

// First returns just a pointer to the first value in the chain.  If the chain is empty, returns `nil`.  Also returns an error if any previous chain method generated an error that affects the first value.
func (receiver *Link[T]) First() (*T, error) {
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

// Last returns just a pointer to the last value in the chain.  If the chain is empty, returns `nil`.  Also returns an error if any previous chain method generated an error that affects the last value.
func (receiver *Link[T]) Last() (*T, error) {
	var lastValue *T
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

// AllMatch will run the `allMatchFunction` parameter function across all the values in the chain.  If every `allMatchFunction` function invocation returns true, this method returns true.  If a single `allMatchFunction` function invocation returns false, this method returns false.  Also returns false and an error if any previous chain method generated an error or if an error is returned from the `allMatchFunction` function.
func (receiver *Link[T]) AllMatch(allMatchFunction func(T) (bool, error)) (bool, error) {
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

// AnyMatch will run the `anyMatchFunction` parameter function across all the values in the chain.  If any `anyMatchFunction` function invocation returns true, this method returns true.  If every invocation `anyMatchFunction` invocation returns false, this method returns false.   Also returns false and an error if any previous chain method generated an error or if an error is returned from the `anyMatchFunction` function.
func (receiver *Link[T]) AnyMatch(anyMatchFunction func(T) (bool, error)) (bool, error) {
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

// NoneMatch will do the exact opposite of `AnyMatch` when it comes to the boolean return value.  Returns an error for the same reasons as `AnyMatch`.
func (receiver *Link[T]) NoneMatch(noneMatchFunction func(T) (bool, error)) (bool, error) {
	match, err := receiver.AnyMatch(noneMatchFunction)
	return !match, err
}

// Reduce applies the `reduceFunction` parameter function to two values in the chain cumulatively.  Subsequent calls to `reduceFunction` uses the previous return value from `reduceFunction` as the first argument and the next value in the chain as the second argument.  A pointer to the final value is returned.  If the chain is empty, `nil` is returned.  Also returns an error if any previous chain method generated an error or if an error is returned from the `reduceFunction` function.
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

// ReduceWithInitialValue applies the `reduceFunction` parameter function to two values in the chain cumulatively.  Subsequent calls to `reduceFunction` uses the previous return value from `reduceFunction` as the first argument and the next value in the chain as the second argument.  The parameter `initialValue` is placed before the entire chain and therefore is the first argument on the first invocation of `initialValue`.  The final value is returned.  If the chain is empty, `initialValue` is returned.  Also returns an error if any previous chain method generated an error or if an error is returned from the `reduceFunction` function.
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
