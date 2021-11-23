// Package keyvalue exists to specify the interface of `KeyValuer`.
package keyvalue

// KeyValuer is used to represent key and value pairs.
type KeyValuer interface {
	Key()   interface{}
	Value() interface{}
}
