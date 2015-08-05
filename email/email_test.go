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
	config.Port = "25"
	config.TLS = false

	var e = NewHtmlMessage("title", "<a href='http://www.google.com'>Google</a>")
	e.To = []string{"*****"}

	fmt.Println(SendMail(config, e))
}