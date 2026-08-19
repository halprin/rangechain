package generator

// mapTuple implements the `keyvalue.KeyValuer` interface and is used to represent map's keys and values.
type mapTuple[K, V any] struct {
	TheKey   K
	TheValue V
}

func (m *mapTuple[K, V]) Key() K {
	return m.TheKey
}

func (m *mapTuple[K, V]) Value() V {
	return m.TheValue
}
