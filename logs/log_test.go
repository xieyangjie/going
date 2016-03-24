package logs

import (
//	"github.com/smartwalle/going/email"
	"testing"
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

//func Test_Mail(t *testing.T) {
//	var config = &email.MailConfig{}
//	config.Username = "*****"
//	config.Host = "*****"
//	config.Password = "*****"
//	config.Port = "25"
//
//	var mailWriter = NewMailWriter(LOG_LEVEL_DEBUG)
//	mailWriter.SetConfig(config)
//	mailWriter.SetToMailList([]string{"917996695@qq.com"})
//
//	SharedLogger().AddOutput("mail", mailWriter)
//
//	writeMessage()
//}

func writeMessage() {
	Debugln("这是Debug消息")
	Infoln("这是Info消息")
	Warnln("这是Warn消息")

	SharedLogger().Flush()
//	Panicln("这是Panic消息")
//	Fatalln("这是Fatal消息")
}