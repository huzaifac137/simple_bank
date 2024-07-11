package util

import (
	"net/smtp"
)

const (
	senderEmail    = "19011598-066@uog.edu.pk"
	senderPassword = "ddwwcrwzqwkisrxx"
)

func SendEmail(to string) error {
	auth := smtp.PlainAuth("", senderEmail, senderPassword, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, senderEmail, []string{to}, []byte("Subject: Verify your email\n\nClick here to verify your email"))
	if err != nil {
		return err
	}
	return nil
}
