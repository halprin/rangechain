package helper

type set struct {
	values map[interface{}]struct{}
}

func NewSet() *set {
	return &set{
		values: make(map[interface{}]struct{}),
	}
}

func (receiver *set) Add(value interface{}) {
	receiver.values[value] = struct{}{}
}

func (receiver *set) Contains(value interface{}) bool {
	_, contains := receiver.values[value]
	return contains
}
