package services

import (
	"clothes-shop-api/internal/config"
	"fmt"
	"net/smtp"
	"strings"
)

type EmailService struct {
	cfg config.Config
}

func NewEmailService(cfg config.Config) *EmailService {
	return &EmailService{cfg: cfg}
}

func (e *EmailService) SendVerificationEmail(toEmail, token, baseURL string) error {
	subject := "Verify Your Email - Clothes Shop"
	verificationLink := fmt.Sprintf("%s/auth/verify-email?token=%s", baseURL, token)

	body := fmt.Sprintf(`
Hello,

Thank you for registering with Clothes Shop!

Please click the link below to verify your email address:

%s

If you did not create an account, please ignore this email.

Best regards,
Clothes Shop Team
`, verificationLink)

	return e.sendEmail(toEmail, subject, body)
}

func (e *EmailService) sendEmail(to, subject, body string) error {
	if e.cfg.SMTPHost == "" || e.cfg.SMTPPort == "" {
		return fmt.Errorf("SMTP configuration is missing")
	}

	auth := smtp.PlainAuth("", e.cfg.SMTPUsername, e.cfg.SMTPPassword, e.cfg.SMTPHost)

	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("From: %s\r\n", e.cfg.EmailFrom))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("\r\n")
	msg.WriteString(body)

	err := smtp.SendMail(
		e.cfg.SMTPHost+":"+e.cfg.SMTPPort,
		auth,
		e.cfg.EmailFrom,
		[]string{to},
		[]byte(msg.String()),
	)

	return err
}
