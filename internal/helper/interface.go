package helper

import "reflect"

func InterfaceSlice(slice interface{}) []interface{} {
	if !IsSlice(slice) {
		panic("non-slice type provided")
	}

	concreteValue := reflect.ValueOf(slice)

	if concreteValue.IsNil() {
		return nil
	}

	interfaceSlice := make([]interface{}, concreteValue.Len())

	for currentIndex := 0; currentIndex < concreteValue.Len(); currentIndex++ {
		interfaceSlice[currentIndex] = concreteValue.Index(currentIndex).Interface()
	}

	return interfaceSlice
}

func InterfaceMap(aMap interface{}) map[interface{}]interface{} {
	if !IsMap(aMap) {
		panic("non-map type provided")
	}

	concreteValue := reflect.ValueOf(aMap)

	if concreteValue.IsNil() {
		return nil
	}

	interfaceMap := make(map[interface{}]interface{}, concreteValue.Len())

	for _, key := range concreteValue.MapKeys() {
		interfaceMap[key.Interface()] = concreteValue.MapIndex(key).Interface()
	}

	return interfaceMap
}
