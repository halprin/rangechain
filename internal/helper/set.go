package helper

type Set[T comparable] struct {
	values map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		values: make(map[T]struct{}),
	}
}

func (receiver *Set[T]) Add(value T) {
	receiver.values[value] = struct{}{}
}

func (receiver *Set[T]) Contains(value T) bool {
	_, contains := receiver.values[value]
	return contains
}
