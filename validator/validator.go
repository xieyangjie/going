package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
	"strconv"
)

const (
	K_VALIDATOR_TAG_NAME          = "validator"
	K_VALIDATOR_TAG_NO_VALIDATION = "-"
)

type ValidatorFunc func(current interface{}, field interface{}, param interface{}) bool

////////////////////////////////////////////////////////////////////////////////
type validatorTag struct {
	Field   string
	Name    string
	Value   interface{}
	Code    int
	Message string
}

func newValidatorTag(field string, name string, value interface{}, code int, message string) *validatorTag {
	var vt = validatorTag{}
	vt.Field = field
	vt.Name  = name
	vt.Value = value
	vt.Code  = code
	vt.Message = message
	return &vt
}

////////////////////////////////////////////////////////////////////////////////
type Validator struct {
	fieldList        map[string]string
	errorList        map[string][]*ValidatorError
	validatorTagList map[string]map[string]*validatorTag
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

func (this *Validator) AddValidator(fieldName string, validatorName string, value interface{}, code int, message string) {
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

	for i := 0; i < numField; i++ {
		var fieldValue = structValue.Field(i)
		var fieldType = structType.Field(i)

		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			fieldValue = fieldValue.Elem()
		}

		this.fieldList[fmt.Sprintln("%d", len(this.fieldList))] = fieldType.Name

		var tag = fieldType.Tag.Get(K_VALIDATOR_TAG_NAME)
		if tag == K_VALIDATOR_TAG_NO_VALIDATION {
			continue
		}

		if tag == "" && ((fieldValue.Kind() != reflect.Struct && fieldValue.Kind() != reflect.Interface) || fieldValue.Type() == reflect.TypeOf(time.Time{})) {
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
						var code, _ = strconv.Atoi(item[2].(string))
						this.AddValidator(name, item[0].(string), item[1], code, item[3].(string))
					}
				}
			}
		} else {
			if len(tagObjList) == 4 {
				var code, _ = strconv.Atoi(tagObjList[2].(string))
				this.AddValidator(name, tagObjList[0].(string), tagObjList[1], code, tagObjList[3].(string))
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
