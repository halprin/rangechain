# rangechain

[![Go Reference](https://pkg.go.dev/badge/github.com/halprin/rangechain.svg)](https://pkg.go.dev/github.com/halprin/rangechain)

Chain together lazily computed modifications to range-able containers.
E.g. slices, arrays, maps, and channels.

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

| Function | Arguments | Description |
| --- | --- | --- |
| `FromSlice` | • `slice` - A slice to be used to start the chain. | Starts the chain with the supplied slice.  Chaining and terminating methods can now be called on the result. |
| `FromArray` | • `array` - An array to be used to start the chain. | Starts the chain with the supplied array.  Chaining and terminating methods can now be called on the result. |
| `FromChannel` | • `channel` - A channel to be used to start the chain. | Starts the chain with the supplied channel.  Chaining and terminating methods can now be called on the result. |
| `FromMap` | • `aMap` - A map to be used to start the chain. | Starts the chain with the supplied map.  Chaining and terminating methods can now be called on the result.  The singular value used to represent the key and value pairs is `keyvalue.KeyValuer` of `github.com/halprin/rangechain/keyvalue`. |

From there, one can call a plethora of additional methods to modify the container passed in originally.  The methods are
outlined below.  The methods fall into one of two categories: chaining or terminating.

#### Value Types

Because Go does not support generics (yet), this library operates using `interface{}`.  One using this library may need
to use type assertions.

For example, something like `actualValue := value.(int)`.

For another example...

```go
stringSlice := []string{"DogCows", "goes", "Moof", "Do", "you", "like", "Clarus", "the", "DogCow"}
chain := FromSlice(stringSlice)

outputSlice, _ := chain.
    Map(func(value interface{}) (interface{}, error) {
        stringValue := value.(string)
        return stringValue + " not", nil
    }).Slice()
```

Notice `stringValue := value.(string)` above.  This allows one to do the string concatenation on the next line because
the `+` operator doesn't work on an `interface{}` type.

### Continuing the Chain

Chaining methods apply some modification to the values in the container values, but keeps the chain alive.
This allows additional chaining methods to be subsequently called on the result.  The subsequent chain methods operate
on any changes performed by the previous chain method.
Because modifications are lazily computed, none of the modifications from chaining methods are applied until _after_ a
terminating method is called.

| Method | Description |
| --- | --- |
| `Map` | Will run the `mapFunction` parameter function parameter against all the values in the chain.  In that function, return what you want to change the value into or an optional error if an error is encountered. |
| `MapParallel` | Will run the `mapFunction` parameter function against all the values in the chain in parallel.  In that function, return what you want to change the value into or an optional error if an error is encountered.  There is overhead to running in parallel so benchmark to ensure you benefit from this version. |
| `Filter` | Will run the `filterFunction` parameter function against all the values in the chain.  In that function, on return of true, the value will stay, or on false, the value will be removed. |
| `FilterParallel` | Will run the `filterFunction` parameter function against all the values in the chain in parallel.  In that function, on return of true, the value will stay in the chain, or on false, the value will be removed from the chain.  There is overhead to running in parallel so benchmark to ensure you benefit from this version. |
| `Skip` | Skips over the parameter `skipNumber` number of values and effectively removes them from the chain.  Also skips over any errors previously generated. |
| `Limit` | Stops the chain after the parameter `keepSize` number of values.  Any elements afterward are effectively removed. |
| `Distinct` | Removes any duplicates. |
| `Flatten` | Will iterate over all the values in the chain, but any value encountered that is a range-able container itself will also have its values iterated over first before continuing with the remaining values in the chain.  Maps flatten to its `keyvalue.KeyValuer` key and value pairs. |
| `Sort` | Sorts the chain given the `Less` function returned from the `returnLessFunction` function parameter.  The `returnLessFunction` function is called with the entire serialized chain as a slice and _returns_ a function that satisfies the same requirements as the [Interface type's](https://pkg.go.dev/sort#Interface) `Less` function.  See the [`TestSortingMaps` example](./example_test.go).  This method is expensive because it must serialize all the values into a slice first. |
| `Reverse` | Reverses the order of the chain.  The last item will be first, and the first item will be last.  This method is expensive because it must serialize all the values into a slice first. |

### Terminating the Chain

Terminating methods also apply some modification, requests some information, or executes something on the values.
They stop the chaining by returning an actual value.  This value will depend on all the previous chaining methods being
executed first.

| Method | Description |
| --- | --- |
| `Slice` | Serializes the chain into a slice and returns it.  Also returns an error if any previous chain method generated an error.  If an error is returned, the slice is filled in until the error was encountered. |
| `Channel` | Serializes the chain into a channel.  Also returns any errors in a channel if any previous chain method generated an error.  If an error is returned, the value channel is closed, the error is sent on the error channel, and the error channel is closed. |
| `ForEach` | Will run the `forEachFunction` parameter function across all the values in the chain.  Also returns an error if any previous chain method generated an error.  If an error is encountered, the function stops executing against the remaining chain. |
| `ForEachParallel` | Will run the `forEachFunction` parameter function across all the values in the chain in parallel.  Also returns an error if any previous chain method generated an error.  If an error is encountered, the function stops executing against the remaining chain.  There is overhead to running in parallel so benchmark to ensure you benefit from this version. |
| `Count` | Returns the number of values in the chain.  Also returns an error if any previous chain method generated an error.  Count returns an accurate number even if an error is encountered. |
| `First` | Returns just a pointer to the first value in the chain.  If the chain is empty, returns `nil`.  Also returns an error if any previous chain method generated an error that affects the first value. |
| `Last` | Returns just a pointer to the last value in the chain.  If the chain is empty, returns `nil`.  Also returns an error if any previous chain method generated an error that affects the last value. |
| `AllMatch` | Will run the `allMatchFunction` parameter function across all the values in the chain.  If every `allMatchFunction` function invocation returns true, this method returns true.  If a single `allMatchFunction` function invocation returns false, this method returns false.  Also returns false and an error if any previous chain method generated an error or if an error is returned from the `allMatchFunction` function. |
| `AnyMatch` | Will run the `anyMatchFunction` parameter function across all the values in the chain.  If any `anyMatchFunction` function invocation returns true, this method returns true.  If every invocation `anyMatchFunction` invocation returns false, this method returns false.   Also returns false and an error if any previous chain method generated an error or if an error is returned from the `anyMatchFunction` function. |
| `NoneMatch` | Will do the exact opposite of `AnyMatch` when it comes to the boolean return value.  Returns an error for the same reasons as `AnyMatch`. |
| `Reduce` | Applies the `reduceFunction` parameter function to two values in the chain cumulatively.  Subsequent calls to `reduceFunction` uses the previous return value from `reduceFunction` as the first argument and the next value in the chain as the second argument.  A pointer to the final value is returned.  If the chain is empty, `nil` is returned.  Also returns an error if any previous chain method generated an error or if an error is returned from the `reduceFunction` function. |
| `ReduceWithInitialValue` | Applies the `reduceFunction` parameter function to two values in the chain cumulatively.  Subsequent calls to `reduceFunction` uses the previous return value from `reduceFunction` as the first argument and the next value in the chain as the second argument.  The parameter `initialValue` is placed before the entire chain and therefore is the first argument on the first invocation of `initialValue`.  The final value is returned.  If the chain is empty, `initialValue` is returned.  Also returns an error if any previous chain method generated an error or if an error is returned from the `reduceFunction` function. |

## Examples

See [`example_test.go` file](./example_test.go) for some examples.
