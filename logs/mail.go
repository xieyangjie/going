package logs

import (
	"fmt"
	"time"
	"github.com/smartwalle/going/email"
)

type MailWriter struct {
	level   int
	config	*email.MailConfig
	to		[]string
}

func NewMailWriter(level int) *MailWriter {
	var writer = &MailWriter{}
	writer.level = level
	return writer
}

func(this *MailWriter) SetLevel(level int) {
	this.level = level
}

func(this *MailWriter) GetLevel() int {
	return this.level
}

func (this *MailWriter) SetConfig(config *email.MailConfig) {
	this.config = config
}

func (this *MailWriter) GetConfig() *email.MailConfig {
	return this.config
}

func (this *MailWriter) SetToMailList(to []string) {
	this.to = to
}

func (this *MailWriter) GetToMailList() []string {
	return this.to
}

func(this *MailWriter) WriteMessage(level int, file string, line int, prefix string, msg string) {
	if level < this.level {
		return
	}

	if this.config == nil {
		return
	}

	if len(this.to) == 0 {
		return
	}

	var message = fmt.Sprintf("%s %s [%s:%d] %s", time.Now().String(), prefix, file, line, msg)
	var mail = email.NewTextMessage(file, message)
	mail.To = this.to

	go email.SendMail(this.config, mail)
}

func(this *MailWriter) Close() {

}

func(this *MailWriter) Flush() {

}