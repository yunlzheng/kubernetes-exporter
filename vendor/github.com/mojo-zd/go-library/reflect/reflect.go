package reflect

import (
	"reflect"
)

func NewInstance(v interface{}) interface{} {
	ty := GetType(reflect.TypeOf(v))
	return reflect.New(ty).Interface()
}

func GetType(ty reflect.Type) (t reflect.Type) {
	if ty.Kind() == reflect.Ptr {
		t = ty.Elem()
		return
	}
	t = ty
	return
}

func GetValue(value reflect.Value) (v reflect.Value) {
	if value.Kind() == reflect.Ptr {
		v = value.Elem()
		return
	}
	v = value
	return
}
