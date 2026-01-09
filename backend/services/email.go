package services

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body string) error {
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	fromName := os.Getenv("MAIL_FROM_NAME")
	fromAddress := os.Getenv("MAIL_FROM_ADDRESS")

	if host == "" || port == "" || username == "" || password == "" {
		return fmt.Errorf("mail configuration is not set")
	}

	auth := smtp.PlainAuth("", username, password, host)

	message := []byte(
		fmt.Sprintf("From: %s <%s>\r\n", fromName, fromAddress) +
			fmt.Sprintf("To: %s\r\n", to) +
			fmt.Sprintf("Subject: %s\r\n", subject) +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n" +
			body,
	)

	return smtp.SendMail(
		host+":"+port,
		auth,
		fromAddress,
		[]string{to},
		message,
	)
}