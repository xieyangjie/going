package tools

import (
	"reflect"
)

func ContainsKeyInMap(m interface{}, key interface{}) bool {
	var mapValue = reflect.ValueOf(m)
	var keyValue = reflect.ValueOf(key)

	if !mapValue.IsValid() || !keyValue.IsValid() {
		return false
	}

	if mapValue.Type().Kind() != reflect.Map {
		return false
	}

	var value = mapValue.MapIndex(keyValue)

	if value.IsValid() {
		return true
	}
	return false
}