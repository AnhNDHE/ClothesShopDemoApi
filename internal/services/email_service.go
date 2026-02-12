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

func (e *EmailService) SendAccountCreatedEmail(toEmail, password string) error {
	subject := "Your Account Has Been Created - Clothes Shop"

	body := fmt.Sprintf(`
Hello,

Your account has been successfully created for Clothes Shop!

Account Details:
- Email: %s
- Password: %s

Please keep this information secure. You can now log in to your account.

If you have any questions, please contact our support team.

Best regards,
Clothes Shop Team
`, toEmail, password)

	return e.sendEmail(toEmail, subject, body)
}

func (e *EmailService) SendEmailUpdateConfirmation(toEmail, token, baseURL string) error {
	subject := "Confirm Your Email Update - Clothes Shop"
	confirmationLink := fmt.Sprintf("%s/auth/confirm-email-update?token=%s", baseURL, token)

	body := fmt.Sprintf(`
Hello,

You have requested to update your email address for your Clothes Shop account.

Please click the link below to confirm this email update:

%s

If you did not request this change, please ignore this email.

Best regards,
Clothes Shop Team
`, confirmationLink)

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
