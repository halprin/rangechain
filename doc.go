// Package rangechain chains together lazily computed modifications to range-able containers.  E.g. slices, arrays, maps, and channels.
//
// To start the chain, use the `From*` functions.  They take a range-able container, and you'll pick the appropriate one for the type of container.
//
//	container := []int{987, 3, 26}
//	chain := rangechain.FromSlice(container)
//
// From there, one can call a plethora of additional methods to modify the container passed in originally.  The methods are
// outlined below.  The methods fall into one of two categories: chaining or terminating.
//
// Value Types
//
// Because Go does not support generics (yet), this library operates using `interface{}`.  One using this library may need
// to use type assertions.
//
// For example, something like `actualValue := value.(int)`.
//
// For another example...
//
//	stringSlice := []string{"DogCows", "goes", "Moof", "Do", "you", "like", "Clarus", "the", "DogCow"}
//	chain := FromSlice(stringSlice)
//
//	outputSlice, _ := chain.
//	    Map(func(value interface{}) (interface{}, error) {
//	        stringValue := value.(string)
//	        return stringValue + " not", nil
//	    }).Slice()
//
// Notice `stringValue := value.(string)` above.  This allows one to do the string concatenation on the next line because
// the `+` operator doesn't work on an `interface{}` type.
//
// Continuing the Chain
//
// Chaining methods apply some modification to the values in the container values, but keeps the chain alive.
// This allows additional chaining methods to be subsequently called on the result.  The subsequent chain methods operate
// on any changes performed by the previous chain method.
// Because modifications are lazily computed, none of the modifications from chaining methods are applied until _after_ a
// terminating method is called.
//
// Terminating the Chain
//
// Terminating methods also apply some modification, requests some information, or executes something on the values.
// They stop the chaining by returning an actual value.  This value will depend on all the previous chaining methods being
// executed first.
package rangechain
