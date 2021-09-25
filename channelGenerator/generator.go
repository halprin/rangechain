package channelGenerator


func FromSlice(slice []interface{}) chan interface{} {
	generation := make(chan interface{})

	go func() {
		for _, currentValue := range slice {
			generation <- currentValue
		}

		close(generation)
	}()

	return generation
}
