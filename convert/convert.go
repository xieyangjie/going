package convert

import (
	"strconv"
	"reflect"
)

func String(value interface{}) string {
	if v, ok := value.(string); ok {
		return v
	}
	return stringValue(value)
}

func Bool(value interface{}) bool {
	if v, ok := value.(bool); ok {
		return v
	}

	var vValue = reflect.ValueOf(value)
	var vKind = vValue.Kind()

	switch vKind {
	case reflect.String:
		var v = vValue.String()
		if v == "true" || v == "yes" || v == "on" || v == "t" || v == "y" || v == "1" {
			return true
		}
		return false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if vValue.Int() == 1 {
			return true
		}
		return false
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if vValue.Uint() == 1 {
			return true
		}
		return false
	case reflect.Float32, reflect.Float64:
		if vValue.Float() > 0.9990 {
			return true
		}
		return false
	case reflect.Bool:
		return vValue.Bool()
	}
	return false
}

func Int(value interface{}) int {
	if v, ok := value.(int); ok {
		return v
	}
	return int(Float64(value))
}

func Int8(value interface{}) int8 {
	if v, ok := value.(int8); ok {
		return v
	}
	return int8(Float64(value))
}

func Int16(value interface{}) int16 {
	if v, ok := value.(int16); ok {
		return v
	}
	return int16(Float64(value))
}

func Int32(value interface{}) int32 {
	if v, ok := value.(int32); ok {
		return v
	}
	return int32(Float64(value))
}

func Int64(value interface{}) int64 {
	if v, ok := value.(int64); ok {
		return v
	}
	return int64(Float64(value))
}

func Uint(value interface{}) uint {
	if v, ok := value.(uint); ok {
		return v
	}
	return uint(Float64(value))
}

func Uint8(value interface{}) uint8 {
	if v, ok := value.(uint8); ok {
		return v
	}
	return uint8(Float64(value))
}

func Uint16(value interface{}) uint16 {
	if v, ok := value.(uint16); ok {
		return v
	}
	return uint16(Float64(value))
}

func Uint32(value interface{}) uint32 {
	if v, ok := value.(uint32); ok {
		return v
	}
	return uint32(Float64(value))
}

func Uint64(value interface{}) uint64 {
	if v, ok := value.(uint64); ok {
		return v
	}
	return uint64(Float64(value))
}

func Float32(value interface{}) float32 {
	if v, ok := value.(float32); ok {
		return v
	}
	return float32(Float64(value))
}

func Float64(value interface{}) float64 {
	if v, ok := value.(float64); ok {
		return v
	}
	return floatValue(value)
}

func floatValue(value interface{}) float64 {
	var vValue = reflect.ValueOf(value)
	var kind = vValue.Kind()

	switch kind {
	case reflect.Bool:
		var v = vValue.Bool()
		if v {
			return 1.0
		}
		return 0.0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(vValue.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(vValue.Int())
	case reflect.Float32, reflect.Float64:
		return vValue.Float()
	case reflect.String:
		var f, err = strconv.ParseFloat(vValue.String(), 64)
		if err == nil {
			return f
		}
	}
	return 0.0
}

func stringValue(value interface{}) string {
	var vValue = reflect.ValueOf(value)
	var vKind =vValue.Kind()

	switch vKind {
	case reflect.Bool:
		return strconv.FormatBool(vValue.Bool())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(vValue.Uint(), 10)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(vValue.Int(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(vValue.Float(), 'f', 6, 32)
	case reflect.String:
		return vValue.String()
	}
	return ""
}
