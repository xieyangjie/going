package validator

import (
	"testing"
	"fmt"
)

type Human struct {
	Name		string		`validator:"[['required', '', '1000', 'name字段不能为空'], ['regex', 'name', '1002', 'name字段只能为英文字母']]"`//`validator:"regex(value:name, code:1001, message:name只能为英文字母)"`
	Age			int			`validator:"['lt', '100', '1004', '老妖精']"`
}

type Student struct {
	Human
}

func Test_Human(t *testing.T) {
	AddRegex("name", "^[a-zA-Z]+$")

	var v = NewValidator()

	var h1 = Student{}
	h1.Name = "q"
	h1.Age = 110

	v.AddValidator("Age", "gte", "18", "1003", "age必须大于等于18")

	fmt.Println(v.Validate(h1), "  ", v.Error())

	fmt.Println(v.Errors())

}