package intermediate

type Link struct {
	generator func() int
}

func NewLink(generator func() int) *Link {
	return &Link{
		generator: generator,
	}
}
