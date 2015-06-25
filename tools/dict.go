package tools

import (
	"reflect"
)

func getValueWithMap(m interface{}, key interface{}) (value reflect.Value) {
	var mapValue = reflect.ValueOf(m)
	var keyValue = reflect.ValueOf(key)

	if !mapValue.IsValid() || !keyValue.IsValid() {
		return
	}
	if mapValue.Type().Kind() != reflect.Map {
		return
	}

	value = mapValue.MapIndex(keyValue)
	return value
}

func ContainsKeyInMap(m interface{}, key interface{}) bool {
	var value = getValueWithMap(m, key)
	if value.IsValid() {
		return true
	}
	return false
}

func GetValueWithMap(m interface{}, key interface{}, defalut interface{}) interface{} {
	var value = getValueWithMap(m, key)
	if value.IsValid() {
		return value.Interface()
	}
	return defalut
}