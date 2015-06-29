package form

import (
	"testing"
	"fmt"
)

type SignInForm struct {
	Form
}

func Test_SignInForm(t *testing.T) {
	var form = NewForm()

	var usernameField = NewField("username", "",
		NewRequired(true, 1000, "用户名不能为空"),
		NewRegex("^[a-z]+[a-z0-9_]{4,39}$", 1001, "用户名格式不正确"),
	)
	form.AddField(usernameField)

	form.AddField(NewField("password", "",  NewRegex("^[\\S]{4,40}$", 1002, "密码格式不正确")))
	form.AddField(NewField("re_password", "", NewEqualToField(form.GetField("password"), 1003, "两次输入的密码不一致")))

	form.BindMap(map[string]string{"username":"test_user", "password":"test", "re_password":"testd"})


	if !form.Validate() {
		var errors map[string][]*FormError = form.GetErrors()

		for _, errs := range errors {
			for _, err := range errs {
				fmt.Println(err.Code, " ", err.Message)
			}
		}
	} else {
		fmt.Println(form.CleanData())
	}
}