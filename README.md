# rangechain
Chain together lazily computed modifications to range-able containers.
E.g. slices, arrays, maps, and channels.

## Install

Require `github.com/halprin/rangechain` in your `go.mod` file or use `go` to add it.
```bash
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
| `FromMap` | • `aMap` - A map to be used to start the chain. | Starts the chain with the supplied map.  Chaining and terminating methods can now be called on the result.  The singular value used to represent the key and value pairs is `generator.MapTuple` of `github.com/halprin/rangechain/generator`. |

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

| Method | Arguments | Description |
| --- | --- | --- |
| `Map` | • `mapFunction` - A function that takes a single argument and returns a value and an optional error. | Will run the function across all the values in the chain.  Return what you want to change the value into and an optional error if an error is encountered. |
| `Filter` | • `filterFunction` - A function that takes a single argument and returns a boolean and an optional error. | Will run the function across all the values in the chain.  On return of true, the value will stay; on false, the value will be removed. |
| `Skip` |  |  |
| `Limit` |  |  |
| `Distinct` |  |  |
| `Flatten` |  |  |
| `Sort` |  |  |
| `Reverse` |  |  |

### Terminating the Chain

Terminating methods also apply some modification, requests some information, or executes something on the elements.
They stop the chaining by returning an actual value.  This value will depend on all the previous chaining methods being
executed first.

## Examples
