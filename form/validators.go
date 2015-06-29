package form

import (
	"regexp"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
type IValidator interface {
	Validate(value string) (bool, *FormError)
}

////////////////////////////////////////////////////////////////////////////////
type Validator struct {
	Code 	int
	Message string
}

func (this *Validator) error() *FormError {
	return NewFormError(this.Code, this.Message)
}

func (this *Validator) Validate(value string) (bool, *FormError) {
	return false, this.error()
}

////////////////////////////////////////////////////////////////////////////////
type Required struct {
	Validator
	Optional bool
}

func NewRequired(optional bool, code int, message string) *Required {
	var required = Required{}
	required.Optional = optional
	required.Code = code
	required.Message = message
	return &required
}

func (this *Required) Validate(value string) (bool, *FormError) {
	if len(strings.TrimSpace(value)) > 0 {
		return true, nil
	}
	return false, this.error()
}

////////////////////////////////////////////////////////////////////////////////
type Regex struct {
	Validator
	Pattern string
}

func NewRegex(pattern string, code int, message string) *Regex {
	var regex = Regex{}
	regex.Pattern = pattern
	regex.Code = code
	regex.Message = message
	return &regex
}

func (this *Regex) Validate(value string) (bool, *FormError) {
	var regex = regexp.MustCompile(this.Pattern)
	var result = regex.MatchString(value)
	if result {
		return true, nil
	}
	return false, this.error()
}

////////////////////////////////////////////////////////////////////////////////
type EqualToValue struct {
	Validator
	Value string
}

func NewEqualToValue(value string, code int, message string) *EqualToValue {
	var equal = EqualToValue{}
	equal.Value = value
	equal.Code = code
	equal.Message = message
	return &equal
}

func (this *EqualToValue) Validate(value string) (bool, *FormError) {
	if value == this.Value {
		return true, nil
	}
	return false, this.error()
}

////////////////////////////////////////////////////////////////////////////////
type EqualToField struct {
	Validator
	Field IField
}

func NewEqualToField(field IField, code int, message string) *EqualToField {
	var equal = EqualToField{}
	equal.Field = field
	equal.Code = code
	equal.Message = message
	return &equal
}

func (this *EqualToField) Validate(value string) (bool, *FormError) {
	if value == this.Field.GetValue() {
		return true, nil
	}
	return false, this.error()
}

////////////////////////////////////////////////////////////////////////////////


////////////////////////////////////////////////////////////////////////////////


////////////////////////////////////////////////////////////////////////////////


////////////////////////////////////////////////////////////////////////////////
