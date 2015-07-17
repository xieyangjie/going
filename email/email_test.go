package email

import (
	"fmt"
	"testing"
)

func Test_SendEmail(t *testing.T) {
	var config = &MailConfig{}
	config.Username = "smartwalle@126.com"
	config.Host = "smtp.126.com"
	config.Password = "yy123456789"
	config.Port = "25"

	var email = NewHtmlEmail("title", "<a href='http://www.baidu.com'>baidu</a>")
	email.To = []string{"917996695@qq.com"}

	fmt.Println(email.Send(config))
}
