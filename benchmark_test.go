package slice_chain

import (
	"testing"
	"time"
)

var size10Slice = makeIntSliceOfSize(10)
var size100Slice = makeIntSliceOfSize(100)
var size1000Slice = makeIntSliceOfSize(1000)
var sliceOfSlice = [][]int{size1000Slice, size1000Slice, size1000Slice}

func BenchmarkIntermediate10(b *testing.B) {
	benchmarkIntermediate(b, size10Slice)
}

func BenchmarkChannelIntermediate10(b *testing.B) {
	benchmarkChannelIntermediate(b, size10Slice)
}

func BenchmarkIntermediate100(b *testing.B) {
	benchmarkIntermediate(b, size100Slice)
}

func BenchmarkChannelIntermediate100(b *testing.B) {
	benchmarkChannelIntermediate(b, size100Slice)
}

func BenchmarkIntermediate1000(b *testing.B) {
	benchmarkIntermediate(b, size1000Slice)
}

func BenchmarkChannelIntermediate1000(b *testing.B) {
	benchmarkChannelIntermediate(b, size1000Slice)
}

func BenchmarkFlatten1000(b *testing.B) {
	for runIndex := 0; runIndex < b.N; runIndex++ {
		FromSlice(sliceOfSlice).Flatten().Filter(func(value interface{}) bool {
			intValue := value.(int)
			return intValue % 2 == 0
		}).Map(func(value interface{}) interface{} {
			intValue := value.(int)
			return intValue * 2 + 2
		}).Slice()
	}
}

func BenchmarkSleepWithSerialMap(b *testing.B) {
	for runIndex := 0; runIndex < b.N; runIndex++ {
		FromSlice(size10Slice).Map(func(value interface{}) interface{} {
			time.Sleep(100 * time.Millisecond)
			return value
		}).Slice()
	}
}

func BenchmarkSleepWithParallelMap(b *testing.B) {
	for runIndex := 0; runIndex < b.N; runIndex++ {
		FromSlice(size10Slice).MapParallel(func(value interface{}) interface{} {
			time.Sleep(100 * time.Millisecond)
			return value
		}).Slice()
	}
}

func makeIntSliceOfSize(size int) []int {
	slice := make([]int, size)

	for index := 0; index < size; index++ {
		slice[index] = index
	}

	return slice
}

func benchmarkIntermediate(b *testing.B, inputSlice []int) {
	for runIndex := 0; runIndex < b.N; runIndex++ {
		FromSlice(inputSlice).Filter(func(value interface{}) bool {
			intValue := value.(int)
			return intValue % 2 == 0
		}).Map(func(value interface{}) interface{} {
			intValue := value.(int)
			return intValue * 2 + 2
		}).Slice()
	}
}

func benchmarkChannelIntermediate(b *testing.B, inputSlice []int) {
	for runIndex := 0; runIndex < b.N; runIndex++ {
		FromSliceWithChannels(inputSlice).Filter(func(value interface{}) bool {
			intValue := value.(int)
			return intValue % 2 == 0
		}).Map(func(value interface{}) interface{} {
			intValue := value.(int)
			return intValue * 2 + 2
		}).Slice()
	}
}
