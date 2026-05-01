package services

import (
	"fmt"
	"net/smtp"
)

type MailService struct {
	smtpHost string
	smtpPort int
	username string
	password string
}

func NewMailService(smtpHost string, smtpPort int, username, password string) *MailService {
	return &MailService{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		username: username,
		password: password,
	}
}

func (s *MailService) SendVerificationEmail(toEmail, code string) error {
	header := "Subject: Հաստատման կոդ - Travel System\n" +
		"MIME-version: 1.0;\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n\n"

	body := fmt.Sprintf(`
        <div style="font-family: sans-serif; border: 1px solid #ddd; padding: 20px; border-radius: 10px; text-align: center;">
            <h2 style="color: #2563eb;">Հաստատման կոդ</h2>
            <div style="font-size: 24px; font-weight: bold; padding: 10px; background: #f3f4f6;">%s</div>
        </div>
    `, code)

	msg := []byte(header + body)
	auth := smtp.PlainAuth("", s.username, s.password, s.smtpHost)
	addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)
	return smtp.SendMail(addr, auth, s.username, []string{toEmail}, msg)
}
