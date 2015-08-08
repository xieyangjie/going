package email

import (
	"fmt"
	"testing"
)

func Test_SendEmail(t *testing.T) {
	var config = &MailConfig{}
	config.Username = "*****"
	config.Host = "*****"
	config.Password = "*****"
	config.Port = "*****"
	config.Secure = true

	var e = NewHtmlMessage("title", "<a href='http://www.google.com'>Google</a>")
	e.To = []string{"917996695@qq.com"}

	fmt.Println(SendMail(config, e))
}