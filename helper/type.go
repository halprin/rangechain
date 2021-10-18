package helper

import "reflect"

func IsSlice(value interface{}) bool {
	concreteValue := reflect.ValueOf(value)
	return concreteValue.Kind() == reflect.Slice
}

func IsArray(value interface{}) bool {
	concreteValue := reflect.ValueOf(value)
	return concreteValue.Kind() == reflect.Array
}
