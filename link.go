package rangechain

// Link is not meant to be initialized directly by external users.  Use the `From*` functions.
type Link[T any] struct {
	generator func() (T, error)
}

func newLink[T any](generator func() (T, error)) *Link[T] {
	return &Link[T]{
		generator: generator,
	}
}
