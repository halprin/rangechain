package intermediate

type Link struct {
	generator func() (interface{}, error)
}

func NewLink(generator func() (interface{}, error)) *Link {
	return &Link{
		generator: generator,
	}
}
