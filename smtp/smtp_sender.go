package smtp

import (
	"fmt"
	"net/smtp"
	"strconv"
)

// MailSender represents an SMTP mail sender.
type MailSender struct {
	smtpServer, sender, pwd, nickname string
	smtpPort                          int
}

// NewMailSender creates a new MailSender instance with the provided SMTP server details.
func NewMailSender(smtpServer string, smtpPort int, sender, pwd, nickname string) *MailSender {
	return &MailSender{
		smtpServer: smtpServer,
		sender:     sender,
		pwd:        pwd,
		nickname:   nickname,
		smtpPort:   smtpPort,
	}
}

// SendMail sends an email using the configured SMTP server.
func (m *MailSender) SendMail(recipient, subject, body string) error {
	auth := smtp.PlainAuth("", m.sender, m.pwd, m.smtpServer)

	message := []byte(fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\nContent-Type: text/html; charset=UTF-8\n\r\n%s\r\n",
		m.nickname, m.sender, recipient, subject, body))

	return smtp.SendMail(m.smtpServer+":"+strconv.Itoa(m.smtpPort), auth, m.sender, []string{recipient}, message)
}
