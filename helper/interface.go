package helper

import "reflect"

func InterfaceSlice(slice interface{}) []interface{} {
	concreteValue := reflect.ValueOf(slice)
	if concreteValue.Kind() != reflect.Slice {
		panic("non-slice type provided")
	}

	if concreteValue.IsNil() {
		return nil
	}

	interfaceSlice := make([]interface{}, concreteValue.Len())

	for currentIndex := 0; currentIndex < concreteValue.Len(); currentIndex++ {
		interfaceSlice[currentIndex] = concreteValue.Index(currentIndex).Interface()
	}

	return interfaceSlice
}
