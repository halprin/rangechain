package rangechain

// Link is not meant to be initialized directly by external users.  Use the `From*` functions.
type Link struct {
	generator func() (interface{}, error)
}

func newLink(generator func() (interface{}, error)) *Link {
	return &Link{
		generator: generator,
	}
}
