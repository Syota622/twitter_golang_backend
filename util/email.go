package util

import (
	"fmt"
	"net/smtp"
)

const (
	SMTPHost = "mailcatcher"
	SMTPPort = 1025
)

func SendEmail(recipient, subject, body string) error {
	from := "no-reply@example.com"
	msg := []byte("To: " + recipient + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", SMTPHost, SMTPPort),
		nil,
		from,
		[]string{recipient},
		msg,
	)

	return err
}
