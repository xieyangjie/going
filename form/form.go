package form

import (
	"sort"
	"net/http"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
type FormFields []IField

func (this FormFields) Len() int {
	return len(this)
}

func (this FormFields) Less(i, j int) bool {
	return this[i].GetIndex() < this[j].GetIndex()
}

func (this FormFields) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

////////////////////////////////////////////////////////////////////////////////
type Form struct {
	formErrors 		map[string][]*FormError
	fields 			map[string]IField
	fieldCounter	int
}

func NewForm() *Form {
	var form = Form{}
	form.formErrors = make(map[string][]*FormError)
	form.fields = make(map[string]IField)
	form.fieldCounter = 0
	return &form
}

func (this *Form) Fields() []IField {
	var fields = make(FormFields, len(this.fields))
	var index = 0
	for _, f := range this.fields {
		fields[index] = f
		index += 1
	}
	sort.Sort(fields)
	return fields
}

func (this *Form) AddField(field IField) {
	if field != nil {
		field.SetIndex(this.fieldCounter)
		this.fields[field.GetName()] = field
		this.fieldCounter += 1
	}
}

func (this *Form) RemoveField(name string) {
	delete(this.fields, name)
}

func (this *Form) GetField(name string) IField {
	return this.fields[name]
}

func (this *Form) ValueWithFieldName(name string) string {
	var field IField = this.GetField(name)
	if field != nil {
		return field.GetValue()
	}
	return ""
}

func (this *Form) BindRequest(request *http.Request) {
	request.ParseForm()
	for key, field := range this.fields {
		var value = request.Form.Get(key)
		if len(strings.TrimSpace(value)) > 0 {
			field.SetValue(value)
		}
	}
}

func (this *Form) BindMap(m map[string]string) {
	for key, value := range m {
		var field = this.GetField(key)
		if field != nil {
			field.SetValue(value)
		}
	}
}

func (this *Form) validate() bool {
	var valid = true
	for _, field := range this.fields {
		var v, e =  field.Validate()
		if !v {
			valid = v
			this.formErrors[field.GetName()] = e
		}
	}
	return valid
}

func (this *Form) Validate() bool {
	return this.validate()
}

func (this *Form) ValidationError() *FormError {
	var fields = this.Fields()
	for _, field := range fields {
		var errs = this.GetErrorWithField(field.GetName())
		if len(errs) > 0 {
			return errs[0]
		}
	}
	return nil
}

func (this *Form) GetErrorWithField(name string) []*FormError {
	return this.formErrors[name]
}

func (this *Form) GetErrors() map[string][]*FormError {
	return this.formErrors
}

func (this *Form) Get(name string) string {
	var field = this.GetField(name)
	if field != nil {
		return field.GetValue()
	}
	return ""
}

func (this *Form) CleanData() map[string]string {
	var data = make(map[string]string)
	for key, field := range this.fields {
		data[key] = field.GetValue()
	}
	return data
}