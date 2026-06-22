// Package rangechain chain together lazily computed modifications to range-able containers. E.g. slices, maps, channels, and iterators.
//
// rangechain uses Go 1.27 generic methods, so no type assertions are needed.
//
// # Start the Chain
//
// To start the chain, use the `From*` functions. They take a range-able container, and you'll pick the appropriate one for the type of container.
//
//	container := []int{987, 3, 26}
//	chain := rangechain.FromSlice(container)
//
// From there, you can call a plethora of additional methods to modify the container passed in originally. The methods fall into one of two categories: chaining or terminating.
//
// # Types
//
// Thanks to Go 1.27 generic methods, the chain's element types flow through every step. No `value.(string)` type assertion needed; `value` is already a `string`.
//
// # Continuing the Chain
//
// Chaining methods apply some modification to the values in the container values but keep the chain alive.
// This allows additional chaining methods to be subsequently called on the result. The subsequent chain methods operate
// on any changes performed by the previous chain method.
// Because modifications are lazily computed, none of the modifications from chaining methods are applied until _after_ a
// terminating method is called.
//
// # Terminating the Chain
//
// Terminating methods also apply some modification, request some information, or execute something on the values.
// They stop the chaining by returning an actual value. This value will depend on all the previous chaining methods being
// executed first.
package rangechain
