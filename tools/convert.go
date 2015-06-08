package tools

import (
	"strconv"
)

func GetString(value interface{}) string {
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

func GetBool(value interface{}) bool {
	if v, ok := value.(bool); ok {
		return v
	}
	v := GetInt(value)
	if v > 0 {
		return true
	}
	return false
}

func GetInt(value interface{}) int {
	if v, ok := value.(int); ok {
		return v
	}
	return int(GetFloat64(value))
}

func GetInt8(value interface{}) int8 {
	if v, ok := value.(int8); ok {
		return v
	}
	return int8(GetFloat64(value))
}

func GetInt16(value interface{}) int16 {
	if v, ok := value.(int16); ok {
		return v
	}
	return int16(GetFloat64(value))
}

func GetInt32(value interface{}) int32 {
	if v, ok := value.(int32); ok {
		return v
	}
	return int32(GetFloat64(value))
}

func GetInt64(value interface{}) int64 {
	if v, ok := value.(int64); ok {
		return v
	}
	return int64(GetFloat64(value))
}

func GetFloat32(value interface{}) float32 {
	if v, ok := value.(float32); ok {
		return v
	}
	return float32(GetFloat64(value))
}

func GetFloat64(value interface{}) float64 {
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