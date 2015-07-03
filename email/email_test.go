package email
import (
	"testing"
	"fmt"
)

func Test_SendEmail(t *testing.T) {
	var config = &MailConfig{}
	config.Username = "邮箱账号"
	config.Host = "smtp.163.com"
	config.Password = "邮箱密码"
	config.Port = "25"

	var email = NewHtmlEmail("title", "<a href='http://www.baidu.com'>baidu</a>")
	email.To = []string{"917996695@qq.com"}

	fmt.Println(email.Send(config))
}