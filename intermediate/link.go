package intermediate

// Link is not meant to be initialized directly by external users.  Use the `From*` functions in the parent package rangechain.
type Link struct {
	generator func() (interface{}, error)
}

// NewLink is not meant to be called directly by external users.  Use the `From*` functions in the parent package rangechain.
func NewLink(generator func() (interface{}, error)) *Link {
	return &Link{
		generator: generator,
	}
}
