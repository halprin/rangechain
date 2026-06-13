// Package rangechain chains together lazily computed modifications to range-able containers.  E.g. slices, maps, channels, and iterators.
//
// rangechain uses Go 1.27 generic methods (GOEXPERIMENT=genericmethods) so chain values keep their static type end-to-end — no type assertions in your callbacks.
//
// To start the chain, use the From* functions.  They take a range-able container, and you'll pick the appropriate one for the type of container.
//
//	container := []int{987, 3, 26}
//	chain := rangechain.FromSlice(container)
//
// From there, one can call a plethora of additional methods to modify the container passed in originally.  The methods are
// outlined below.  The methods fall into one of two categories: chaining or terminating.
//
// # Value Types
//
// Thanks to Go 1.27 generic methods, the chain's element type flows through every step.  Map can change the element type;
// Filter, Skip, Limit, etc. preserve it.
//
//	stringSlice := []string{"DogCows", "goes", "Moof", "Do", "you", "like", "Clarus", "the", "DogCow"}
//
//	outputSlice, _ := rangechain.FromSlice(stringSlice).
//	    Map(func(value string) (string, error) {
//	        return value + " not", nil
//	    }).Slice()
//
// No value.(string) cast needed — value is already a string, and outputSlice is a []string.
//
// # Continuing the Chain
//
// Chaining methods apply some modification to the values in the container values, but keeps the chain alive.
// This allows additional chaining methods to be subsequently called on the result.  The subsequent chain methods operate
// on any changes performed by the previous chain method.
// Because modifications are lazily computed, none of the modifications from chaining methods are applied until _after_ a
// terminating method is called.
//
// # Terminating the Chain
//
// Terminating methods also apply some modification, requests some information, or executes something on the values.
// They stop the chaining by returning an actual value.  This value will depend on all the previous chaining methods being
// executed first.
//
// # Flatten
//
// Flatten is a chainable generic method.  The caller specifies the output element type U; each emitted inner value is type-asserted to U, and a mismatch injects an error into the chain.
//
//	result, _ := rangechain.FromSlice([][]int{{1, 2, 3}, {4, 5, 6}}).
//	    Flatten[int]().Slice()
//
// Use Flatten[any]() for heterogeneous chains that mix containers and scalars.
package rangechain
