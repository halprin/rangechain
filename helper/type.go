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

func IsChannel(value interface{}) bool {
	concreteValue := reflect.ValueOf(value)
	return concreteValue.Kind() == reflect.Chan
}

func IsMap(value interface{}) bool {
	concreteValue := reflect.ValueOf(value)
	return concreteValue.Kind() == reflect.Map
}
