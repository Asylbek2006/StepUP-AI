package email

import (
	"fmt"
	"net/smtp"
)

type SMTPEmailSender struct {
	host     string
	port     string
	username string
	password string
}

func NewSMTPEmailSender(host, port, username, password string) *SMTPEmailSender {
	return &SMTPEmailSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (s *SMTPEmailSender) SendPasswordResetEmail(toEmail, resetLink string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	subject := "Subject: StepUp AI - Password Reset\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	body := fmt.Sprintf(`
		<h2>Password Reset</h2>
		<p>Click the link below to reset your password:</p>
		<a href="%s">Reset Password</a>
		<p>This link expires in 1 hour.</p>
	`, resetLink)

	message := []byte(subject + mime + "\r\n" + body)

	address := fmt.Sprintf("%s:%s", s.host, s.port)
	return smtp.SendMail(address, auth, s.username, []string{toEmail}, message)
}
