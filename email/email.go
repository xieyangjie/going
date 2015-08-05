package email

import (
	"errors"
	"net/mail"
	"net/smtp"
	"crypto/tls"
	"net"
)

////////////////////////////////////////////////////////////////////////////////
type MailConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	TLS		 bool
}

func (this *MailConfig) Address() string {
	return this.Host + ":" + this.Port
}

////////////////////////////////////////////////////////////////////////////////
// Send an email using the given host and SMTP auth (optional), returns any error thrown by smtp.SendMail
// This function merges the To, Cc, and Bcc fields and calls the smtp.SendMail function using the Email.Bytes() output as the message
func send(addr string, a smtp.Auth, m *Message, tls bool) error {
	// Merge the To, Cc, and Bcc fields
	to := make([]string, 0, len(m.To)+len(m.Cc)+len(m.Bcc))
	to = append(append(append(to, m.To...), m.Cc...), m.Bcc...)
	for i := 0; i < len(to); i++ {
		addr, err := mail.ParseAddress(to[i])
		if err != nil {
			return err
		}
		to[i] = addr.Address
	}
	// Check to make sure there is at least one recipient and one "From" address
	if m.from == "" || len(to) == 0 {
		return errors.New("Must specify at least one From address and one To address")
	}
	from, err := mail.ParseAddress(m.from)
	if err != nil {
		return err
	}
	raw, err := m.Bytes()
	if err != nil {
		return err
	}

	if tls {
		return SendTLSMail(addr, a, from.Address, to, raw)
	}
	return smtp.SendMail(addr, a, from.Address, to, raw)
}

func SendMail(config *MailConfig, m *Message) error {
	if config == nil {
		return errors.New("config 不能为空")
	}
	if m.from != config.Username {
		m.from = config.Username
	}

	return send(config.Address(), smtp.PlainAuth("", config.Username, config.Password, config.Host), m, config.TLS)
}

func SendTLSMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {

	host, _, _ := net.SplitHostPort(addr)

	tlsConfig := &tls.Config {
		InsecureSkipVerify: true,
		ServerName: host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}