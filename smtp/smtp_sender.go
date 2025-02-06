package smtp

import (
	"gopkg.in/gomail.v2"
)

// MailSender represents an SMTP mail sender.
type MailSender struct {
	sender, nickname string
	cli              *gomail.Dialer
}

// NewMailSender creates a new MailSender instance with the provided SMTP server details.
func NewMailSender(smtpServer string, smtpPort int, sender, pwd, nickname string) *MailSender {
	dialer := gomail.NewDialer(smtpServer, smtpPort, sender, pwd)
	return &MailSender{
		sender:   sender,
		nickname: nickname,
		cli:      dialer,
	}
}

// SendMail sends an email using the configured SMTP server.
func (m *MailSender) SendMail(recipient, subject, body, bodyContentType string, cc []string) error {
	msg := gomail.NewMessage(gomail.SetCharset("UTF-8"))
	msg.SetAddressHeader("From", m.sender, m.nickname)
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", subject)
	msg.SetBody(bodyContentType, body)
	if len(cc) > 0 {
		msg.SetHeader("Cc", cc...)
	}

	return m.cli.DialAndSend(msg)
}
