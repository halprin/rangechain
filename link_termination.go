package rangechain

import (
	"errors"
	"iter"

	"github.com/halprin/rangechain/internal/generator"
)

// Slice serializes the chain into a slice and returns it. Also returns an error if any previous chain method generated one. On error, the slice is filled in until the error was encountered.
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

// Channel serializes the chain into a channel. Returns a paired error channel. If an error occurs, the value channel is closed, the error is sent on the error channel, then the error channel is closed.
func (receiver *Link[T]) Channel() (<-chan T, <-chan error) {
	endChannel := make(chan T)
	errorChannel := make(chan error)

	go func() {
		for {
			currentValue, err := receiver.generator()
			if err != nil {
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

// Iterator returns an `iter.Seq2[T, error]` so the chain can be consumed with `range`. Yields `(value, nil)` for each value; if an upstream error occurs, yields `(zero, err)` once and stops.
func (receiver *Link[T]) Iterator() iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for {
			currentValue, err := receiver.generator()
			if err != nil {
				if errors.Is(err, generator.Exhausted) {
					return
				}
				var zero T
				yield(zero, err)
				return
			}

			if !yield(currentValue, nil) {
				return
			}
		}
	}
}

// ForEach runs `forEachFunction` parameter across every value in the chain. Stops on the first error and returns it.
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

// ForEachParallel is like ForEach, but invocations run concurrently. There is overhead to running in parallel so benchmark to ensure you benefit from this version.
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

// Count returns the number of values in the chain. Counts accurately even when an error occurs partway, and returns that error.
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

// First returns a pointer to the first value. `nil` if the chain is empty. Also returns an error if encountered.
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

// Last returns a pointer to the last value. `nil` if the chain is empty. Also returns an error if encountered.
func (receiver *Link[T]) Last() (*T, error) {
	var lastValue *T
	var lastError error

	for {
		currentValue, err := receiver.generator()
		if err != nil && errors.Is(err, generator.Exhausted) {
			return lastValue, lastError
		}

		lastValue = new(currentValue)
		lastError = err
	}
}

// AllMatch returns true when the `allMatchFunction` parameter returns true for every value, false otherwise. False (with the error) when an error is encountered or if `allMatchFunction` errors itself.
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

// AnyMatch returns true when the `anyMatchFunction` parameter returns true for any value, false otherwise. False (with the error) when an error is encountered or if `anyMatchFunction` errors itself.
func (receiver *Link[T]) AnyMatch(anyMatchFunction func(T) (bool, error)) (bool, error) {
	for {
		currentValue, err := receiver.generator()
		if err != nil {
			if errors.Is(err, generator.Exhausted) {
				return false, nil
			} else if !errors.Is(err, generator.Exhausted) {
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

// NoneMatch is the boolean opposite of `AnyMatch`. Returns an error for the same reasons as `AnyMatch`.
func (receiver *Link[T]) NoneMatch(noneMatchFunction func(T) (bool, error)) (bool, error) {
	match, err := receiver.AnyMatch(noneMatchFunction)
	return !match, err
}

// Reduce runs the `reduceFunction` parameter to two values in the chain cumulatively. Subsequent calls to `reduceFunction` uses the previous return value from `reduceFunction` as the first argument and the next value in the chain as the second argument. A pointer to the final value is returned. If the chain is empty, `nil` is returned. Also returns an error if any previous chain method generated an error or if an error is returned from the `reduceFunction` function.
func (receiver *Link[T]) Reduce(reduceFunction func(T, T) (T, error)) (*T, error) {
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
		err = nil
	}

	return &intermediateItem, err
}

// ReduceWithInitialValue is similar to Reduce, but starts with `initialValue` in the chain.
func (receiver *Link[T]) ReduceWithInitialValue[A any](reduceFunction func(A, T) (A, error), initialValue A) (A, error) {
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
		err = nil
	}

	return intermediateItem, err
}
