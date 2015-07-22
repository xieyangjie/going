package convert

import (
	"strconv"
	"reflect"
)

func ConvertToString(value interface{}) string {
	if v, ok := value.(string); ok {
		return v
	}
	return strconv.FormatFloat(ConvertToFloat64(value), 'f', -1, 64)
}

func ConvertToBool(value interface{}) bool {
	if v, ok := value.(bool); ok {
		return v
	}
	v := ConvertToInt(value)
	if v > 0 {
		return true
	}
	return false
}

func ConvertToInt(value interface{}) int {
	if v, ok := value.(int); ok {
		return v
	}
	return int(ConvertToFloat64(value))
}

func ConvertToInt8(value interface{}) int8 {
	if v, ok := value.(int8); ok {
		return v
	}
	return int8(ConvertToFloat64(value))
}

func ConvertToInt16(value interface{}) int16 {
	if v, ok := value.(int16); ok {
		return v
	}
	return int16(ConvertToFloat64(value))
}

func ConvertToInt32(value interface{}) int32 {
	if v, ok := value.(int32); ok {
		return v
	}
	return int32(ConvertToFloat64(value))
}

func ConvertToInt64(value interface{}) int64 {
	if v, ok := value.(int64); ok {
		return v
	}
	return int64(ConvertToFloat64(value))
}

func ConvertToUint(value interface{}) uint64 {
	if v, ok := value.(uint64); ok {
		return v
	}
	return uint64(ConvertToFloat64(value))
}

func ConvertToFloat32(value interface{}) float32 {
	if v, ok := value.(float32); ok {
		return v
	}
	return float32(ConvertToFloat64(value))
}

func ConvertToFloat64(value interface{}) float64 {
	if v, ok := value.(float64); ok {
		return v
	}
	return floatValue(value)
}

func floatValue(value interface{}) float64 {
	var vf = reflect.ValueOf(value)
	var k = vf.Kind()

	switch k {
	case reflect.Bool:
		 var v = vf.Bool()
		if v {
			return 1.0
		}
		return 0.0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(vf.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(vf.Int())
	case reflect.Float32, reflect.Float64:
		return vf.Float()
	case reflect.String:
		var f, err = strconv.ParseFloat(vf.String(), 64)
		if err == nil {
			return f
		}
	}
	return 0.0
}
