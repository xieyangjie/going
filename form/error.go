package form

type FormError struct {
	Code 	int
	Message string
}

func NewFormError(code int, message string) *FormError {
	var formError = &FormError{}
	formError.Code = code
	formError.Message = message
	return formError
}