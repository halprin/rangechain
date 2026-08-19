# rangechain

[![Go Reference](https://pkg.go.dev/badge/github.com/halprin/rangechain.svg)](https://pkg.go.dev/github.com/halprin/rangechain)

Chain together lazily computed modifications to range-able containers.
E.g. slices, maps, channels, and iterators.

`rangechain` uses Go 1.27 generic methods, so no type assertions are needed.

## Install

Require `github.com/halprin/rangechain` in your `go.mod` file or use `go` to add it.
```shell
$ go get github.com/halprin/rangechain
```

## API

Start by importing the package.
```go
import "github.com/halprin/rangechain"
```

### Examples

Filter then map:

```go
words := []string{"DogCow", "Moof", "Clarus", "the", "DogCow"}

result, _ := rangechain.FromSlice(words).
    Filter(func(w string) (bool, error) { return len(w) > 3, nil }).
    Map(func(w string) (string, error) { return strings.ToUpper(w), nil }).
    Slice()
// result == []string{"DOGCOW", "MOOF", "CLARUS", "DOGCOW"}
```

`Map` can change the element type:

```go
lengths, _ := rangechain.FromSlice([]string{"a", "bb", "ccc"}).
    Map(func(w string) (int, error) { return len(w), nil }).
    Slice()
// lengths == []int{1, 2, 3}
```

Sort a map by value, then take the top N keys:

```go
scores := map[string]int{"Mac OS 9": 9, "Mac OS X": 10, "System 7": 7}

top, _ := rangechain.FromMap(scores).
    Sort(func(kvs []keyvalue.KeyValuer[string, int]) func(int, int) bool {
        return func(i, j int) bool { return kvs[i].Value() > kvs[j].Value() }
    }).
    Map(func(kv keyvalue.KeyValuer[string, int]) (string, error) {
        return kv.Key(), nil
    }).
    Limit(2).
    Slice()
// top == []string{"Mac OS X", "Mac OS 9"}
```

Consume the chain lazily with `Iterator` and `range`:

```go
chain := rangechain.FromSlice([]int{1, 2, 3, 4, 5}).
    Filter(func(v int) (bool, error) { return v%2 == 1, nil }).
    Map(func(v int) (int, error) { return v * 10, nil })

for value, err := range chain.Iterator() {
    if err != nil {
        // handle the upstream error and stop
        break
    }
    fmt.Println(value) // 10, 30, 50
}
```

Reduce into a different type to build a frequency map:

```go
words := []string{"go", "rust", "go", "c"}

counts, _ := rangechain.FromSlice(words).
    ReduceWithInitialValue(func(acc map[string]int, w string) (map[string]int, error) {
        acc[w]++
        return acc, nil
    }, map[string]int{})
// counts == map[string]int{"go": 2, "rust": 1, "c": 1}
```

Homogeneous case — caller knows the inner element type:

```go
result, _ := rangechain.FromSlice([][]int{{1, 2, 3}, {4, 5, 6}}).
    Flatten[int]().Slice()
// result == []int{1, 2, 3, 4, 5, 6}
```

Heterogeneous case — mix containers and scalars, use `Flatten[any]()`:

```go
result, _ := rangechain.FromSlice([]any{
    []int{1, 2, 3}, 4, []int{5, 6},
}).Flatten[any]().Slice()
// result == []any{1, 2, 3, 4, 5, 6}
```

See [`example_test.go` file](./example_test.go) for more examples.

### Start the Chain

To start the chain, use the `From*` functions. They take a range-able container, and you'll pick the appropriate one
for the type of container.

```go
container := []int{987, 3, 26}
chain := rangechain.FromSlice(container)
```

| Function       | Arguments                                                    | Description                                                                                                                                                                                                 |
|----------------|--------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `FromSlice`    | • `slice []T` - A slice to start the chain.                  | Starts the chain with the supplied slice. Pass `arr[:]` for an array. Chaining and terminating methods can now be called on the result.                                                                     |
| `FromChannel`  | • `channel <-chan T` - A channel to start the chain.         | Starts the chain with the supplied channel. Chaining and terminating methods can now be called on the result.                                                                                               |
| `FromMap`      | • `aMap map[K]V` - A map to start the chain.                 | Starts the chain with the supplied map. Each chain element is a `keyvalue.KeyValuer[K, V]` from `github.com/halprin/rangechain/keyvalue`. Chaining and terminating methods can now be called on the result. |
| `FromIterator` | • `anIterator iter.Seq[T]` - An iterator to start the chain. | Starts the chain with the supplied iterator. Chaining and terminating methods can now be called on the result.                                                                                              |

From there, you can call a plethora of additional methods to modify the container passed in originally. The methods fall into one of two categories: chaining or terminating.

#### Types

Thanks to Go 1.27 generic methods, the chain's element types flow through every step. No `value.(string)` type assertion needed; `value` is already a `string`.

### Continuing the Chain

Chaining methods apply some modification to the values in the container values but keep the chain alive.
This allows additional chaining methods to be subsequently called on the result. The subsequent chain methods operate
on any changes performed by the previous chain method.
Because modifications are lazily computed, none of the modifications from chaining methods are applied until _after_ a
terminating method is called.

| Method           | Description                                                                                                                                                                                                                                                                                                                             |
|------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `Map`            | Runs the `mapFunction` parameter against all the values in the chain. In that function, return what you want to change the value into or an optional error.                                                                                                                                                                             |
| `MapParallel`    | Like `Map`, but invocations run concurrently. There is overhead to running in parallel so benchmark to ensure you benefit from this version.                                                                                                                                                                                            |
| `Filter`         | Runs the `filterFunction` parameter against all the values in the chain. Returning true keeps the value; returning false drops it.                                                                                                                                                                                                      |
| `FilterParallel` | Like `Filter`, but invocations run concurrently. There is overhead to running in parallel so benchmark to ensure you benefit from this version.                                                                                                                                                                                         |
| `Skip`           | Skips the next `skipNumber` values in the chain (including any errors already in flight).                                                                                                                                                                                                                                               |
| `Limit`          | Stops the chain after `keepSize` values have been emitted.                                                                                                                                                                                                                                                                              |
| `DistinctFunc`   | Removes duplicates. Two values whose `keyFunction` returns the same value are considered equal.  Use `func(v T) T { return v }` when the values are already comparable.                                                                                                                                                                 |
| `Flatten`        | Iterates each chain value; any value that is itself a slice, channel, iterator, or map is descended into (maps emit `keyvalue.KeyValuer[any, any]` entries). Each emitted inner value is type-asserted to `U`; a mismatch injects an error into the chain at that point.                                                                |
| `Sort`           | Sorts the chain using a `Less` function returned by the `returnLessFunction` parameter. The returned function must satisfy the same requirements as the [Interface type's](https://pkg.go.dev/sort#Interface) `Less` function. See the [`TestSortingMaps` example](./example_test.go). Expensive because it serializes the chain first. |
| `Reverse`        | Reverses the order of the chain. Expensive because it serializes the chain first.                                                                                                                                                                                                                                                       |

### Terminating the Chain

Terminating methods also apply some modification, request some information, or execute something on the values.
They stop the chaining by returning an actual value. This value will depend on all the previous chaining methods being
executed first.

| Method                   | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
|--------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `Slice`                  | Serializes the chain into a slice and returns it. Also returns an error if any previous chain method generated one. On error, the slice is filled in until the error was encountered.                                                                                                                                                                                                                                                                                         |
| `Channel`                | Serializes the chain into a channel. Returns a paired error channel. If an error occurs, the value channel is closed, the error is sent on the error channel, then the error channel is closed.                                                                                                                                                                                                                                                                               |
| `Iterator`               | Returns an `iter.Seq2[T, error]` so the chain can be consumed with `range`. Yields `(value, nil)` for each value; if an upstream error occurs, yields `(zero, err)` once and stops.                                                                                                                                                                                                                                                                                           |
| `ForEach`                | Runs `forEachFunction` parameter across every value in the chain. Stops on the first error and returns it.                                                                                                                                                                                                                                                                                                                                                                    |
| `ForEachParallel`        | Like `ForEach`, but invocations run concurrently. There is overhead to running in parallel so benchmark to ensure you benefit from this version.                                                                                                                                                                                                                                                                                                                              |
| `Count`                  | Returns the number of values in the chain. Counts accurately even when an error occurs partway, and returns that error.                                                                                                                                                                                                                                                                                                                                                       |
| `First`                  | Returns a pointer to the first value. `nil` if the chain is empty. Also returns an error if encountered.                                                                                                                                                                                                                                                                                                                                                                      |
| `Last`                   | Returns a pointer to the last value. `nil` if the chain is empty. Also returns an error if encountered.                                                                                                                                                                                                                                                                                                                                                                       |
| `AllMatch`               | True when the `allMatchFunction` parameter returns true for every value, false otherwise. False (with the error) when an error is encountered or if `allMatchFunction` errors itself.                                                                                                                                                                                                                                                                                         |
| `AnyMatch`               | True when the `anyMatchFunction` parameter returns true for any value, false otherwise. False (with the error) when an error is encountered or if `anyMatchFunction` errors itself.                                                                                                                                                                                                                                                                                           |
| `NoneMatch`              | Boolean opposite of `AnyMatch`. Returns an error for the same reasons as `AnyMatch`.                                                                                                                                                                                                                                                                                                                                                                                          |
| `Reduce`                 | Runs the `reduceFunction` parameter to two values in the chain cumulatively. Subsequent calls to `reduceFunction` uses the previous return value from `reduceFunction` as the first argument and the next value in the chain as the second argument. A pointer to the final value is returned. If the chain is empty, `nil` is returned. Also returns an error if any previous chain method generated an error or if an error is returned from the `reduceFunction` function. |
| `ReduceWithInitialValue` | Similar to `Reduce`, but starts with `initialValue` in the chain.                                                                                                                                                                                                                                                                                                                                                                                                             |
