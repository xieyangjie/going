package logs

import (
	"github.com/smartwalle/going/email"
	"testing"
	"time"
)

func Test_Init(t *testing.T) {
	SharedLogger().SetEnableStack(true)
	SharedLogger().SetStackLevel(LOG_LEVEL_PANIC)
}

func Test_File(t *testing.T) {
	var fileWriter = NewFileWriter(LOG_LEVEL_DEBUG, "./test_logs")
	SharedLogger().AddOutput("file", fileWriter)

	writeMessage()
}

func Test_Mail(t *testing.T) {
	var config = &email.MailConfig{}
	config.Username = "*****"
	config.Host = "*****"
	config.Password = "*****"
	config.Port = "25"

	var mailWriter = NewMailWriter(LOG_LEVEL_DEBUG)
	mailWriter.SetConfig(config)
	mailWriter.SetToMailList([]string{"917996695@qq.com"})

	SharedLogger().AddOutput("mail", mailWriter)

	writeMessage()

	//等待邮件发出，在正常的服务器环境中，不需要这样
	time.Sleep(5 * time.Second)
}

func writeMessage() {
	Debugln("这是Debug消息")
	Infoln("这是Info消息")
	Warnln("这是Warn消息")
	Panicln("这是Panic消息")
	Fatalln("这是Fatal消息")
}