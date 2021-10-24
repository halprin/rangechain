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

### Continuing the Chain

### Terminating the Chain

## Examples
