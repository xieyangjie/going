package form

import "fmt"

////////////////////////////////////////////////////////////////////////////////
type IField interface {
	GetName()	string

	GetValue()	string
	SetValue(value string)

	//用作排序
	GetIndex()	int
	SetIndex(value int)

	Validate() (bool, []*FormError)
}

////////////////////////////////////////////////////////////////////////////////
type Field struct {
	Name 		 	string
	Value 		 	string
	Index		 	int
	Validators 	 	[]IValidator
}

func NewField(name string, defaultValue string, validators ...IValidator) *Field {
	var field = &Field{}
	field.Name = name
	field.Value = defaultValue
	field.Validators = validators
	return field
}

func (this *Field) GetName() string {
	return this.Name
}

func (this *Field) GetValue() string {
	return this.Value
}

func (this *Field) SetValue(value string) {
	this.Value = value
}

func (this *Field) GetIndex() int{
	return this.Index
}

func (this *Field) SetIndex(value int) {
	this.Index = value
}

func (this *Field) String() string {
	return fmt.Sprintf("%s=%s", this.Name, this.Value)
}

func (this *Field) Validate() (bool, []*FormError) {
	var valid = true
	var errs = make([]*FormError, 0, len(this.Validators))
	for _, validator := range this.Validators {
		var v, e = validator.Validate(this.GetValue())
		if !v {
			valid = false
			errs = append(errs, e)
		}
	}
	return valid, errs
}

