package validator

import (
	"github.com/smartwalle/going/convert"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var validatorFuncList = map[string]ValidatorFunc{}

var customerFuncList = map[string]ValidatorFunc{}

var customerRegexList = map[string]*regexp.Regexp{}

////////////////////////////////////////////////////////////////////////////////
func init() {
	registerFunction("required", required)
	registerFunction("len", length)
	registerFunction("eq", equal)
	registerFunction("ne", unequal)
	registerFunction("in", inList)
	registerFunction("nin", notInList)
	registerFunction("lt", lessThan)
	registerFunction("lte", lessThanOrEqual)
	registerFunction("gt", greaterThan)
	registerFunction("gte", greaterThanOrEqual)
	registerFunction("eqf", equalToField)
	registerFunction("nef", unequalToField)
	registerFunction("ltf", lessThanField)
	registerFunction("ltef", lessThanOrEqualToField)
	registerFunction("gtf", greaterThanField)
	registerFunction("gtef", greaterThanOrEqualToField)
	registerFunction("regex", regex)
}

////////////////////////////////////////////////////////////////////////////////
func registerFunction(key string, fun ValidatorFunc) {
	validatorFuncList[key] = fun
}

////////////////////////////////////////////////////////////////////////////////
func AddFunction(key string, fun ValidatorFunc) bool {
	if len(key) == 0 && fun == nil {
		customerFuncList[key] = fun
		return true
	}
	return false
}

func RemoveFunction(key string) {
	delete(customerFuncList, key)
}

////////////////////////////////////////////////////////////////////////////////
func AddRegex(key string, pattern string) bool {
	if len(key) != 0 && len(pattern) != 0 {
		var regex = regexp.MustCompile(pattern)
		if regex != nil {
			customerRegexList[key] = regex
			return true
		}
	}
	return false
}

func RevmoeRegex(key string) {
	if len(key) == 0 {
		return
	}
	delete(customerRegexList, key)
}

////////////////////////////////////////////////////////////////////////////////
func required(current interface{}, field interface{}, param interface{}) bool {
	var fieldValue = reflect.ValueOf(field)

	switch fieldValue.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return field != nil && fieldValue.Len() > 0
	default:
		return field != nil && field != reflect.Zero(reflect.TypeOf(field)).Interface()
	}
}

func length(current interface{}, field interface{}, param interface{}) bool {
	var fieldValue = reflect.ValueOf(field)

	switch fieldValue.Kind() {
	case reflect.String:
		l := convert.ConvertToInt(param)
		return utf8.RuneCountInString(fieldValue.String()) == l
	case reflect.Slice, reflect.Array, reflect.Map:
		l := convert.ConvertToInt(param)
		return fieldValue.Len() == l
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		l := convert.ConvertToInt64(param)
		return fieldValue.Int() == l
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		l := convert.ConvertToUint(param)
		return fieldValue.Uint() == l
	case reflect.Float32, reflect.Float64:
		l := convert.ConvertToFloat64(param)
		return fieldValue.Float() == l
	default:
		return false
	}
}

func equal(current interface{}, field interface{}, param interface{}) bool {
	var fieldValue = reflect.ValueOf(field)

	switch fieldValue.Kind() {
	case reflect.String:
		return strings.EqualFold(fieldValue.String(), param.(string))
	case reflect.Slice, reflect.Array, reflect.Map:
		p := convert.ConvertToInt(param)
		return fieldValue.Len() == p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := convert.ConvertToInt64(param)
		return fieldValue.Int() == p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := convert.ConvertToUint(param)
		return fieldValue.Uint() == p
	case reflect.Float32, reflect.Float64:
		p := convert.ConvertToFloat64(param)
		return fieldValue.Float() == p
	case reflect.Struct:
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
			var t1 = field.(time.Time).Unix()
			var t2 = convert.ConvertToInt64(param)
			return t1 == t2
		}
		return false
	default:
		return false
	}
}

func unequal(current interface{}, field interface{}, param interface{}) bool {
	return !equal(current, field, param)
}

func lessThan(current interface{}, field interface{}, param interface{}) bool {
	var fieldValue = reflect.ValueOf(field)

	switch fieldValue.Kind() {
	case reflect.String:
		return fieldValue.String() < param.(string)
	case reflect.Slice, reflect.Array, reflect.Map:
		p := convert.ConvertToInt(param)
		return fieldValue.Len() < p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := convert.ConvertToInt64(param)
		return fieldValue.Int() < p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := convert.ConvertToUint(param)
		return fieldValue.Uint() < p
	case reflect.Float32, reflect.Float64:
		p := convert.ConvertToFloat64(param)
		return fieldValue.Float() < p
	case reflect.Struct:
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
			var t1 = field.(time.Time).Unix()
			var t2 = convert.ConvertToInt64(param)
			return t1 < t2
		}
		return false
	default:
		return false
	}
}

func inList(current interface{}, field interface{}, param interface{}) bool {
	var fieldValue = reflect.ValueOf(field)

	switch fieldValue.Kind() {
	case reflect.Struct:
	case reflect.Slice, reflect.Array, reflect.Map:
		return false
	default:
		if list, ok := param.([]interface{}); ok {
			for _, item := range list {

				switch fieldValue.Kind() {
				case reflect.String:
					if item.(string) == fieldValue.String() {
						return true
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if convert.ConvertToInt64(item) == fieldValue.Int() {
						return true
					}
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
					if convert.ConvertToUint(item) == fieldValue.Uint() {
						return true
					}
				case reflect.Float32, reflect.Float64:
					if convert.ConvertToFloat64(item) == fieldValue.Float() {
						return true
					}
				}
			}
		}
		return false
	}
	return false
}

func notInList(current interface{}, field interface{}, param interface{}) bool {
	return !inList(current, field, param)
}

func lessThanOrEqual(current interface{}, field interface{}, param interface{}) bool {
	var fieldValue = reflect.ValueOf(field)

	switch fieldValue.Kind() {
	case reflect.String:
		return fieldValue.String() <= param.(string)
	case reflect.Slice, reflect.Array, reflect.Map:
		p := convert.ConvertToInt(param)
		return fieldValue.Len() <= p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := convert.ConvertToInt64(param)
		return fieldValue.Int() <= p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := convert.ConvertToUint(param)
		return fieldValue.Uint() <= p
	case reflect.Float32, reflect.Float64:
		p := convert.ConvertToFloat64(param)
		return fieldValue.Float() <= p
	case reflect.Struct:
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
			var t1 = field.(time.Time).Unix()
			var t2 = convert.ConvertToInt64(param)
			return t1 <= t2
		}
		return false
	default:
		return false
	}
}

func greaterThan(current interface{}, field interface{}, param interface{}) bool {
	var fieldValue = reflect.ValueOf(field)

	switch fieldValue.Kind() {
	case reflect.String:
		return fieldValue.String() > param.(string)
	case reflect.Slice, reflect.Array, reflect.Map:
		p := convert.ConvertToInt(param)
		return fieldValue.Len() > p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := convert.ConvertToInt64(param)
		return fieldValue.Int() > p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := convert.ConvertToUint(param)
		return fieldValue.Uint() > p
	case reflect.Float32, reflect.Float64:
		p := convert.ConvertToFloat64(param)
		return fieldValue.Float() > p
	case reflect.Struct:
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
			var t1 = field.(time.Time).Unix()
			var t2 = convert.ConvertToInt64(param)
			return t1 > t2
		}
		return false
	default:
		return false
	}
}

func greaterThanOrEqual(current interface{}, field interface{}, param interface{}) bool {
	var fieldValue = reflect.ValueOf(field)

	switch fieldValue.Kind() {
	case reflect.String:
		return fieldValue.String() >= param.(string)
	case reflect.Slice, reflect.Array, reflect.Map:
		p := convert.ConvertToInt(param)
		return fieldValue.Len() >= p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := convert.ConvertToInt64(param)
		return fieldValue.Int() >= p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := convert.ConvertToUint(param)
		return fieldValue.Uint() >= p
	case reflect.Float32, reflect.Float64:
		p := convert.ConvertToFloat64(param)
		return fieldValue.Float() >= p
	case reflect.Struct:
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
			var t1 = field.(time.Time).Unix()
			var t2 = convert.ConvertToInt64(param)
			return t1 >= t2
		}
		return false
	default:
		return false
	}
}

func equalToField(current interface{}, field interface{}, param interface{}) bool {
	if current == nil {
		return false
	}
	var structValue = reflect.ValueOf(current)
	if structValue.Kind() == reflect.Ptr && !structValue.IsNil() {
		structValue = reflect.ValueOf(structValue.Elem().Interface())
	}

	var toFieldValue reflect.Value

	switch structValue.Kind() {
	case reflect.Struct:
		if structValue.Type() == reflect.TypeOf(time.Time{}) {
			toFieldValue = structValue
		} else {
			var tempField = structValue.FieldByName(param.(string))
			if tempField.Kind() == reflect.Invalid {
				return false
			}
			toFieldValue = tempField
		}
	default:
		toFieldValue = structValue
	}

	if toFieldValue.Kind() == reflect.Ptr && !toFieldValue.IsNil() {
		toFieldValue = reflect.ValueOf(toFieldValue.Elem().Interface())
	}

	var fieldValue = reflect.ValueOf(field)

	switch fieldValue.Kind() {
	case reflect.String:
		return fieldValue.String() == toFieldValue.String()
	case reflect.Slice, reflect.Array, reflect.Map:
		p := toFieldValue.Len()
		return fieldValue.Len() == p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fieldValue.Int() == toFieldValue.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fieldValue.Uint() == toFieldValue.Uint()
	case reflect.Float32, reflect.Float64:
		return fieldValue.Float() == toFieldValue.Float()
	case reflect.Struct:
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) && toFieldValue.Type() == reflect.TypeOf(time.Time{}) {
			var t1 = fieldValue.Interface().(time.Time)
			var t2 = toFieldValue.Interface().(time.Time)
			return t1.Equal(t2)
		}
		return false
	default:
		return false
	}
}

func unequalToField(current interface{}, field interface{}, param interface{}) bool {
	return !equalToField(current, field, param)
}

func lessThanField(current interface{}, field interface{}, param interface{}) bool {
	if current == nil {
		return false
	}
	var structValue = reflect.ValueOf(current)
	if structValue.Kind() == reflect.Ptr && !structValue.IsNil() {
		structValue = reflect.ValueOf(structValue.Elem().Interface())
	}

	var toFieldValue reflect.Value

	switch structValue.Kind() {
	case reflect.Struct:
		if structValue.Type() == reflect.TypeOf(time.Time{}) {
			toFieldValue = structValue
		} else {
			var tempField = structValue.FieldByName(param.(string))
			if tempField.Kind() == reflect.Invalid {
				return false
			}
			toFieldValue = tempField
		}
	default:
		toFieldValue = structValue
	}

	if toFieldValue.Kind() == reflect.Ptr && !toFieldValue.IsNil() {
		toFieldValue = reflect.ValueOf(toFieldValue.Elem().Interface())
	}

	var fieldValue = reflect.ValueOf(field)

	switch fieldValue.Kind() {
	case reflect.String:
		return fieldValue.String() < toFieldValue.String()
	case reflect.Slice, reflect.Array, reflect.Map:
		p := toFieldValue.Len()
		return fieldValue.Len() < p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fieldValue.Int() < toFieldValue.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fieldValue.Uint() < toFieldValue.Uint()
	case reflect.Float32, reflect.Float64:
		return fieldValue.Float() < toFieldValue.Float()
	case reflect.Struct:
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) && toFieldValue.Type() == reflect.TypeOf(time.Time{}) {
			var t1 = fieldValue.Interface().(time.Time)
			var t2 = toFieldValue.Interface().(time.Time)
			return t1.Before(t2)
		}
		return false
	default:
		return false
	}
}

func lessThanOrEqualToField(current interface{}, field interface{}, param interface{}) bool {
	if current == nil {
		return false
	}
	var structValue = reflect.ValueOf(current)
	if structValue.Kind() == reflect.Ptr && !structValue.IsNil() {
		structValue = reflect.ValueOf(structValue.Elem().Interface())
	}

	var toFieldValue reflect.Value

	switch structValue.Kind() {
	case reflect.Struct:
		if structValue.Type() == reflect.TypeOf(time.Time{}) {
			toFieldValue = structValue
		} else {
			var tempField = structValue.FieldByName(param.(string))
			if tempField.Kind() == reflect.Invalid {
				return false
			}
			toFieldValue = tempField
		}
	default:
		toFieldValue = structValue
	}

	if toFieldValue.Kind() == reflect.Ptr && !toFieldValue.IsNil() {
		toFieldValue = reflect.ValueOf(toFieldValue.Elem().Interface())
	}

	var fieldValue = reflect.ValueOf(field)

	switch fieldValue.Kind() {
	case reflect.String:
		return fieldValue.String() <= toFieldValue.String()
	case reflect.Slice, reflect.Array, reflect.Map:
		p := toFieldValue.Len()
		return fieldValue.Len() <= p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fieldValue.Int() <= toFieldValue.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fieldValue.Uint() <= toFieldValue.Uint()
	case reflect.Float32, reflect.Float64:
		return fieldValue.Float() <= toFieldValue.Float()
	case reflect.Struct:
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) && toFieldValue.Type() == reflect.TypeOf(time.Time{}) {
			var t1 = fieldValue.Interface().(time.Time)
			var t2 = toFieldValue.Interface().(time.Time)
			return (t1.Before(t2) || t1.Equal(t2))
		}
		return false
	default:
		return false
	}
}

func greaterThanField(current interface{}, field interface{}, param interface{}) bool {
	if current == nil {
		return false
	}
	var structValue = reflect.ValueOf(current)
	if structValue.Kind() == reflect.Ptr && !structValue.IsNil() {
		structValue = reflect.ValueOf(structValue.Elem().Interface())
	}

	var toFieldValue reflect.Value

	switch structValue.Kind() {
	case reflect.Struct:
		if structValue.Type() == reflect.TypeOf(time.Time{}) {
			toFieldValue = structValue
		} else {
			var tempField = structValue.FieldByName(param.(string))
			if tempField.Kind() == reflect.Invalid {
				return false
			}
			toFieldValue = tempField
		}
	default:
		toFieldValue = structValue
	}

	if toFieldValue.Kind() == reflect.Ptr && !toFieldValue.IsNil() {
		toFieldValue = reflect.ValueOf(toFieldValue.Elem().Interface())
	}

	var fieldValue = reflect.ValueOf(field)

	switch fieldValue.Kind() {
	case reflect.String:
		return fieldValue.String() > toFieldValue.String()
	case reflect.Slice, reflect.Array, reflect.Map:
		p := toFieldValue.Len()
		return fieldValue.Len() > p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fieldValue.Int() > toFieldValue.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fieldValue.Uint() > toFieldValue.Uint()
	case reflect.Float32, reflect.Float64:
		return fieldValue.Float() > toFieldValue.Float()
	case reflect.Struct:
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) && toFieldValue.Type() == reflect.TypeOf(time.Time{}) {
			var t1 = fieldValue.Interface().(time.Time)
			var t2 = toFieldValue.Interface().(time.Time)
			return t1.After(t2)
		}
		return false
	default:
		return false
	}
}

func greaterThanOrEqualToField(current interface{}, field interface{}, param interface{}) bool {
	if current == nil {
		return false
	}
	var structValue = reflect.ValueOf(current)
	if structValue.Kind() == reflect.Ptr && !structValue.IsNil() {
		structValue = reflect.ValueOf(structValue.Elem().Interface())
	}

	var toFieldValue reflect.Value

	switch structValue.Kind() {
	case reflect.Struct:
		if structValue.Type() == reflect.TypeOf(time.Time{}) {
			toFieldValue = structValue
		} else {
			var tempField = structValue.FieldByName(param.(string))
			if tempField.Kind() == reflect.Invalid {
				return false
			}
			toFieldValue = tempField
		}
	default:
		toFieldValue = structValue
	}

	if toFieldValue.Kind() == reflect.Ptr && !toFieldValue.IsNil() {
		toFieldValue = reflect.ValueOf(toFieldValue.Elem().Interface())
	}

	var fieldValue = reflect.ValueOf(field)

	switch fieldValue.Kind() {
	case reflect.String:
		return fieldValue.String() >= toFieldValue.String()
	case reflect.Slice, reflect.Array, reflect.Map:
		p := toFieldValue.Len()
		return fieldValue.Len() >= p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fieldValue.Int() >= toFieldValue.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fieldValue.Uint() >= toFieldValue.Uint()
	case reflect.Float32, reflect.Float64:
		return fieldValue.Float() >= toFieldValue.Float()
	case reflect.Struct:
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) && toFieldValue.Type() == reflect.TypeOf(time.Time{}) {
			var t1 = fieldValue.Interface().(time.Time)
			var t2 = toFieldValue.Interface().(time.Time)
			return (t1.After(t2) || t1.Equal(t2))
		}
		return false
	default:
		return false
	}
}

func matchesRegex(regex *regexp.Regexp, field interface{}) bool {
	return regex.MatchString(field.(string))
}

func regex(current interface{}, field interface{}, param interface{}) bool {
	var regex *regexp.Regexp
	var ok = false
	if regex, ok = customerRegexList[param.(string)]; !ok {
		regex = regexp.MustCompile(param.(string))
	}

	if regex != nil {
		return matchesRegex(regex, field)
	}
	return false
}
