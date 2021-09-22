package helper

import "reflect"

func IsSlice(value interface{}) bool {
	concreteValue := reflect.ValueOf(value)
	return concreteValue.Kind() == reflect.Slice
}
