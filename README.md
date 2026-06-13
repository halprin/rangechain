# rangechain

[![Go Reference](https://pkg.go.dev/badge/github.com/halprin/rangechain.svg)](https://pkg.go.dev/github.com/halprin/rangechain)

Chain together lazily computed modifications to range-able containers.
E.g. slices, maps, channels, and iterators.

`rangechain` uses Go 1.27 generic methods (`GOEXPERIMENT=genericmethods`) so chain values keep their static type end-to-end — no type assertions in your callbacks.

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

### Start the Chain

To start the chain, use the `From*` functions.  They take a range-able container, and you'll pick the appropriate one
for the type of container.

```go
container := []int{987, 3, 26}
chain := rangechain.FromSlice(container)
```

| Function       | Arguments                                                        | Description                                                                                                                                                                                                                                   |
|----------------|------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `FromSlice`    | • `slice []T` - A slice to start the chain.                      | Starts the chain with the supplied slice.  Pass `arr[:]` for an array.  Returns `*Link[T]`.                                                                                                                                                  |
| `FromChannel`  | • `channel <-chan T` - A channel to start the chain.             | Starts the chain with the supplied channel.  Returns `*Link[T]`.                                                                                                                                                                              |
| `FromMap`      | • `aMap map[K]V` - A map to start the chain.                     | Starts the chain with the supplied map.  Each chain element is a `keyvalue.KeyValuer[K, V]` from `github.com/halprin/rangechain/keyvalue`.  Returns `*Link[keyvalue.KeyValuer[K, V]]`.                                                        |
| `FromIterator` | • `anIterator iter.Seq[T]` - An iterator to start the chain.     | Starts the chain with the supplied iterator.  Returns `*Link[T]`.                                                                                                                                                                             |

From there, you can call a plethora of additional methods to modify the container passed in originally. The methods fall into one of two categories: chaining or terminating.

#### Value Types

Thanks to Go 1.27 generic methods, the chain's element type flows through every step.  `Map` can change the element type; `Filter`, `Skip`, `Limit`, etc. preserve it.

```go
stringSlice := []string{"DogCows", "goes", "Moof", "Do", "you", "like", "Clarus", "the", "DogCow"}

outputSlice, _ := rangechain.FromSlice(stringSlice).
    Map(func(value string) (string, error) {
        return value + " not", nil
    }).Slice()
```

No `value.(string)` cast needed — `value` is already a `string`, and `outputSlice` is a `[]string`.

### Continuing the Chain

Chaining methods apply some modification to the values in the container values, but keeps the chain alive.
This allows additional chaining methods to be subsequently called on the result.  The subsequent chain methods operate
on any changes performed by the previous chain method.
Because modifications are lazily computed, none of the modifications from chaining methods are applied until _after_ a
terminating method is called.

| Method | Description |
| --- | --- |
| `Map[U any]` | Runs `mapFunction` against every value, returning a chain whose element type is `U`. Return the new value, or an optional error. |
| `MapParallel[U any]` | Like `Map`, but invocations run concurrently.  There is overhead to running in parallel — benchmark to confirm you benefit. |
| `Filter` | Runs `filterFunction` against every value.  Returning true keeps the value; returning false drops it. |
| `FilterParallel` | Like `Filter`, but invocations run concurrently.  There is overhead to running in parallel — benchmark to confirm you benefit. |
| `Skip` | Skips the next `skipNumber` values (including any errors already in flight). |
| `Limit` | Stops the chain after `keepSize` values have been emitted. |
| `DistinctFunc[K comparable]` | Removes duplicates.  Two values whose `keyFunction` returns the same `K` are considered equal.  Use `func(v T) T { return v }` when `T` is itself comparable. |
| `Flatten[U any]` | Iterates each chain value; any value that is itself a slice, array, channel, or map is descended into (maps emit `keyvalue.KeyValuer[any, any]` entries).  Each emitted inner value is type-asserted to `U`; a mismatch injects an error into the chain at that point. |
| `Sort` | Sorts the chain using a `Less` function returned by `returnLessFunction`.  The returned function operates against the serialized `[]T`.  Expensive because it serializes the chain first. |
| `Reverse` | Reverses the order of the chain.  Expensive because it serializes the chain first. |

#### Flatten examples

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

### Terminating the Chain

Terminating methods also apply some modification, request some information, or execute something on the values.
They stop the chaining by returning an actual value.  This value will depend on all the previous chaining methods being
executed first.

| Method | Description |
| --- | --- |
| `Slice` | Serializes the chain into a `[]T` and returns it.  Also returns an error if any previous chain method generated one.  On error, the slice is filled in until the error was encountered. |
| `Channel` | Serializes the chain into a `<-chan T`.  Returns a paired `<-chan error`.  If an error occurs, the value channel is closed, the error is sent on the error channel, then the error channel is closed. |
| `ForEach` | Runs `forEachFunction` across every value.  Stops on the first error. |
| `ForEachParallel` | Like `ForEach`, but invocations run concurrently.  Benchmark to confirm benefit. |
| `Count` | Returns the number of values in the chain.  Counts accurately even when an error occurs partway. |
| `First` | Returns a pointer to the first value.  `nil` if the chain is empty. |
| `Last` | Returns a pointer to the last value.  `nil` if the chain is empty. |
| `AllMatch` | True when `allMatchFunction` returns true for every value.  False (with the error) if `allMatchFunction` errors. |
| `AnyMatch` | True when `anyMatchFunction` returns true for any value. |
| `NoneMatch` | Boolean opposite of `AnyMatch`. |
| `Reduce` | Folds the chain with `reduceFunction(prev T, next T) (T, error)`.  Returns `*T`; `nil` if the chain is empty. |
| `ReduceWithInitialValue[A any]` | Folds with `reduceFunction(acc A, next T) (A, error)` starting from `initialValue A`.  The accumulator type can differ from the element type. |

## Examples

See [`example_test.go` file](./example_test.go) for some examples.
