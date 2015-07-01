package validator

import (
	"fmt"
)

type ValidatorError struct {
	Field 		string
	Code		int
	Message		string
	Index		int
}

func NewValidatorError(field string, code int, message string) *ValidatorError {
	var err = ValidatorError{}
	err.Field = field
	err.Code = code
	err.Message = message
	return &err
}

func (this *ValidatorError) Error() string {
	return fmt.Sprintf("[%s]%d:%s", this.Field, this.Code, this.Message)
}