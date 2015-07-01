package validator

import (
	"reflect"
	"time"
	"strings"
//	"regexp"
	"strconv"
	"fmt"
	"encoding/json"
)

const (
	K_VALIDATOR_TAG_NAME 			= "validator"
//	K_VALIDATOR_TAG_SEPARATOR 		= "&"
	K_VALIDATOR_TAG_NO_VALIDATION 	= "-"
//	K_VALIDATOR_TAG_PARAM_SEPARATOR	= ","
//	K_VALIDATOR_TAG_VALUE_SEPARATOR	= ":"

//	K_VALIDATOR_NAME_PATTERN		= `([a-zA-Z]*?)\(`
//	K_VALIDATOR_PARAM_PATTERN		= `\((.*?)\)`
)

////////////////////////////////////////////////////////////////////////////////
type validatorTag struct {
	Field string
	Name string
	Value string
	Code int
	Message string
}

func newValidatorTag(field string, name string, value string, code string, message string) *validatorTag {
	var vt = validatorTag{}
	vt.Field = field
	vt.Name = name
	vt.Value = value
	vt.Code, _ = strconv.Atoi(code)
	vt.Message = message
	return &vt
}

////////////////////////////////////////////////////////////////////////////////
//var (
//	validatorNameRegex = regexp.MustCompile(K_VALIDATOR_NAME_PATTERN)
//	validatorParamRegex = regexp.MustCompile(K_VALIDATOR_PARAM_PATTERN)
//)

type ValidatorFunc func(current interface{}, field interface{}, param string) bool

type Validator struct {
	fieldList 			map[string]string
	errorList 			map[string][]*ValidatorError
	validatorTagList 	map[string]map[string]*validatorTag
}

func NewValidator() *Validator {
	var validator = Validator{}
	validator.fieldList = make(map[string]string)
	validator.errorList = make(map[string][]*ValidatorError)
	validator.validatorTagList = make(map[string]map[string]*validatorTag)
	return &validator
}

func (this *Validator) Error() *ValidatorError {
	var l = len(this.fieldList)
	for i := 0; i < l; i++ {
		var key = fmt.Sprintln("%d", i)
		var errList = this.ErrorsWithField(this.fieldList[key])
		if len(errList) > 0 {
			return errList[0]
		}
	}
	return nil
}

func (this *Validator) Errors() map[string][]*ValidatorError {
	return this.errorList
}

func (this *Validator) ErrorsWithField(name string) []*ValidatorError {
	return this.errorList[name]
}

func (this *Validator) AddValidator(fieldName string, validatorName string, value string, code string, message string) {
	var vt = newValidatorTag(fieldName, validatorName, value, code, message)
	var vtm = this.validatorTagList[fieldName]
	if vtm == nil {
		vtm = make(map[string]*validatorTag)
	}
	vtm[validatorName] = vt
	this.validatorTagList[fieldName] = vtm
}

func (this *Validator) Validate(obj interface{}) bool {
	this.validateStruct(obj, obj)

	if len(this.errorList) == 0 {
		return true
	}
	return false
}

func (this *Validator) validateStruct(current interface{}, s interface{}) {
	var structValue = reflect.ValueOf(s)

	if structValue.Kind() == reflect.Ptr && !structValue.IsNil() {
		this.validateStruct(current, structValue.Elem().Interface())
		return
	}

	var structType = reflect.TypeOf(s)
	var numField = structValue.NumField()

	for i:=0; i<numField; i++ {
		var fieldValue = structValue.Field(i)
		var fieldType = structType.Field(i)

		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			fieldValue = fieldValue.Elem()
		}

		var tag = fieldType.Tag.Get(K_VALIDATOR_TAG_NAME)
		if tag == K_VALIDATOR_TAG_NO_VALIDATION {
			continue
		}

		if tag == "" && ((fieldValue.Kind() != reflect.Struct && fieldValue.Kind() != reflect.Interface ) || fieldValue.Type() == reflect.TypeOf(time.Time{})){
			this.validateField(current, fieldValue.Interface(), fieldType.Name, "")
			continue
		}

		switch fieldValue.Kind() {
		case reflect.Struct, reflect.Interface:
			if !fieldType.Anonymous {
				continue
			}
			if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
				this.validateField(current, fieldValue.Interface(), fieldType.Name, tag)
			} else {
				this.validateStruct(fieldValue.Interface(), fieldValue.Interface())
			}
		default:
			this.fieldList[fmt.Sprintln("%d", len(this.fieldList))] = fieldType.Name
			this.validateField(current, fieldValue.Interface(), fieldType.Name, tag)
		}
	}
}



func (this *Validator) validateField(current interface{}, field interface{}, name string, tagValue string) {
	tagValue = strings.Replace(tagValue, "'", "\"", -1)
	var tagObjList []interface{}
	json.Unmarshal([]byte(tagValue), &tagObjList)

	if len(tagObjList) > 0 {
		var obj = tagObjList[0]
		if _, ok := obj.([]interface{}); ok {
			for _, tagObj := range tagObjList {
				if item, ok := tagObj.([]interface{}); ok {
					if len(item) == 4 {
						this.AddValidator(name, item[0].(string), item[1].(string), item[2].(string), item[3].(string))
					}
				}
			}
		} else {
			if len(tagObjList) == 4 {
				this.AddValidator(name, tagObjList[0].(string), tagObjList[1].(string), tagObjList[2].(string), tagObjList[3].(string))
			}
		}
	}

	var validatorList = this.validatorTagList[name]

	var errList = make([]*ValidatorError, 0, len(validatorList))

	for _, vt := range validatorList {
		var _, err = this.validate(current, field, vt)
		if err != nil {
			errList = append(errList, err)
		}
	}

	if len(errList) > 0 {
		this.errorList[name] = errList
	}
}

//func (this *Validator) validateField(current interface{}, field interface{}, name string, tagValue string) {
//
//	var tagList = strings.Split(tagValue, K_VALIDATOR_TAG_SEPARATOR)
//	var validatorList = this.customerValidator[name]
//
//	var errList = make([]*ValidatorError, 0, len(tagList)+len(validatorList))
//
//	for _, tag := range tagList {
//		var _, err = this.validateTag(current, field, name, strings.TrimSpace(tag))
//		if err != nil {
//			errList = append(errList, err)
//		}
//	}
//
//	for _, vt := range validatorList {
//		var _, err = this.validate(current, field, vt)
//		if err != nil {
//			errList = append(errList, err)
//		}
//	}
//
//	if len(errList) > 0 {
//		this.errorList[name] = errList
//	}
//}
//
//func (this *Validator) getValidatorNameString(tag string) string {
//	var results = validatorNameRegex.FindStringSubmatch(tag)
//	if len(results) > 1 {
//		return strings.TrimSpace(results[1])
//	}
//	return ""
//}
//
//func (this *Validator) getParamList(tag string) map[string]string {
//	var paramStrList = validatorParamRegex.FindStringSubmatch(tag)
//
//	if len(paramStrList) > 1 {
//		var items = strings.Split(paramStrList[1], K_VALIDATOR_TAG_PARAM_SEPARATOR)
//		var result = make(map[string]string)
//		for _, item := range items {
//			var param = strings.Split(item, K_VALIDATOR_TAG_VALUE_SEPARATOR)
//			if len(param) > 1 {
//				result[strings.TrimSpace(param[0])] = strings.TrimSpace(param[1])
//			}
//		}
//		return result
//	}
//	return nil
//}
//
//func (this *Validator) validateTag(current interface{}, field interface{}, name string, tag string) (bool, *ValidatorError) {
//	var validatorName = this.getValidatorNameString(tag)
//	var paramList = this.getParamList(tag)
//
//	var param = paramList["value"]
//	var code = paramList["code"]
//	var message = paramList["message"]
//
//	var vt = newValidatorTag(name, validatorName, param, code, message)
//
//	return this.validate(current, field, vt)
//}

func (this *Validator) validate(current interface{}, field interface{}, tag *validatorTag) (bool, *ValidatorError) {
	valFunc, ok := validatorFuncList[tag.Name]
	if !ok {
		valFunc, _ = customerFuncList[tag.Name]
	}

	if valFunc != nil {
		var result = valFunc(current, field, tag.Value)
		var err *ValidatorError = nil
		if !result {
			err = NewValidatorError(tag.Field, tag.Code, tag.Message)
		}
		return result, err
	}
	return false, nil
}