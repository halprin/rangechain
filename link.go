package rangechain

// Link is not meant to be initialized directly by external users.  Use the `From*` functions in the parent package rangechain.
type Link struct {
	generator func() (interface{}, error)
}

// newLink is not meant to be called directly by external users.  Use the `From*` functions in the parent package rangechain.
func newLink(generator func() (interface{}, error)) *Link {
	return &Link{
		generator: generator,
	}
}
