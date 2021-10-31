package keyvalue

// KeyValuer is used to represent the key and value pairs.
type KeyValuer interface {
	Key()   interface{}
	Value() interface{}
}
