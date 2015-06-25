package convert

import (
	"strconv"
)

func ConvertToString(value interface{}) string {
	if v, ok := value.(string); ok {
		return v
	} else if v, ok := value.(uint64); ok {
		return strconv.FormatUint(v, 10)
	} else if v, ok := value.(int64); ok {
		return strconv.FormatInt(v, 10)
	} else if v, ok := value.(float64); ok {
		return strconv.FormatFloat(v, 'f', -1, 64)
	} else if v, ok := value.(bool); ok {
		return strconv.FormatBool(v)
	}
	return ""
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

func ConvertToFloat32(value interface{}) float32 {
	if v, ok := value.(float32); ok {
		return v
	}
	return float32(ConvertToFloat64(value))
}

func ConvertToFloat64(value interface{}) float64 {
	if v, ok := value.(float64); ok {
		return v
	} else if v, ok := value.(string); ok {
		vf, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return vf
		}
	}
	return 0.0
}