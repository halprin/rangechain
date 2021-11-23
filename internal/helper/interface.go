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
