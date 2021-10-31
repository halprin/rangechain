// Package intermediate is used continue or terminate the chain.
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
package intermediate
