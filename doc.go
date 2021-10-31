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
package rangechain
