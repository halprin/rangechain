package generator

// mapTuple implements the `keyvalue.KeyValuer` interface and is used when
type mapTuple struct {
	TheKey   interface{}
	TheValue interface{}
}

func (m *mapTuple) Key() interface{} {
	return m.TheKey
}

func (m *mapTuple) Value() interface{} {
	return m.TheValue
}
