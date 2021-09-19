package iterator

func FromSlice(theSlice []int) func() int {
	currentIndex := 0
	return func() int {
		value := theSlice[currentIndex]
		currentIndex++
		return value
	}
}
