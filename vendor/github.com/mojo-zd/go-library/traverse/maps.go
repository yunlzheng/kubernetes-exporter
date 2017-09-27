package traverse

import (
	"reflect"

	r "github.com/mojo-zd/go-library/reflect"
)

type MapHandleFunc func(value interface{}) interface{}
type MapIteratorHandleFun func(key, value interface{}, index int)

func MapIterator(m interface{}, iteratorHandleFunc MapIteratorHandleFun) {
	if !isMap(m) {
		panic("m must be map!")
	}

	value := reflect.ValueOf(m)
	for index, k := range value.MapKeys() {
		iteratorHandleFunc(k.Interface(), value.MapIndex(k).Interface(), index)
	}
}

//key只支持基本类型
func ContainsKey(m interface{}, key interface{}) (contains bool) {
	if !isMap(m) {
		panic("m must be map!")
	}

	value := reflect.ValueOf(m)
	for _, k := range value.MapKeys() {
		if k.Interface() == key {
			contains = true
			break
		}
	}
	return
}

func ContainsValue(m interface{}, value interface{}) (contains bool) {
	v := reflect.ValueOf(m)
	if !isMap(m) {
		panic("m must be map!")
	}

	for _, k := range v.MapKeys() {
		if compare(v.MapIndex(k).Interface(), value) {
			contains = true
			break
		}
	}
	return
}

//指定struct中Key的value作为map的key map的value可以由MapHandleFunc的返回值决定, 如果想以struct作为value回调函数设置为nil即可
//ex: Person{
// Name string
// Sex int
// }
// structs := []Person{{Name: "mojo", Sex:1}}
// StructsToMap(structs, "Name")
func StructsToMap(slice interface{}, key string, handleFunc MapHandleFunc) (m interface{}) {
	if !isSlice(slice) {
		panic("collection must be slice!")
		return
	}
	result := map[interface{}]interface{}{}
	v := reflect.ValueOf(slice)
	for index := 0; index < v.Len(); index++ {
		value := r.GetValue(v).Index(index).Interface()
		keyValue := GetValueByName(value, key)

		if handleFunc != nil {
			result[keyValue] = handleFunc(value)
		} else {
			result[keyValue] = value
		}
	}
	m = result
	return
}

func GetValueByName(i interface{}, key string) (value interface{}) {
	ty := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	numField := r.GetType(ty).NumField()

	for index := 0; index < numField; index++ {
		field := r.GetType(ty).Field(index)
		fieldValue := r.GetValue(v)
		if field.Name == key {
			value = fieldValue.FieldByName(field.Name).Interface()
			break
		}
	}
	return
}
