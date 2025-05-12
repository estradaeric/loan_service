package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"time"

	"loan-service/logger"

	"github.com/sirupsen/logrus"
)

// SendEmail sends a basic plain text email (without attachment)
func SendEmail(to, subject, body string) error {
	return SendEmailWithAttachment(to, subject, body, nil, "")
}

// SendEmailWithAttachment sends an email with optional PDF attachment
func SendEmailWithAttachment(to, subject, body string, attachment []byte, filename string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		logger.Log.Error("SMTP environment configuration is incomplete")
		return fmt.Errorf("SMTP configuration is incomplete")
	}

	var msg bytes.Buffer
	boundary := fmt.Sprintf("BOUNDARY-%d", time.Now().UnixNano())

	// Headers
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")

	if attachment != nil {
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
		msg.WriteString("\r\n")

		// Message part
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		msg.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
		msg.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
		msg.WriteString(body + "\r\n")

		// Attachment part (base64 encoded)
		encoded := make([]byte, base64.StdEncoding.EncodedLen(len(attachment)))
		base64.StdEncoding.Encode(encoded, attachment)

		msg.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
		msg.WriteString(fmt.Sprintf("Content-Type: application/pdf; name=\"%s\"\r\n", filename))
		msg.WriteString("Content-Transfer-Encoding: base64\r\n")
		msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", filename))

		for i := 0; i < len(encoded); i += 76 {
			end := i + 76
			if end > len(encoded) {
				end = len(encoded)
			}
			msg.Write(encoded[i:end])
			msg.WriteString("\r\n")
		}
		msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else {
		msg.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
		msg.WriteString(body + "\r\n")
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	if err := smtp.SendMail(addr, auth, from, []string{to}, msg.Bytes()); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"to":     to,
			"subject": subject,
			"error":  err.Error(),
		}).Error("Failed to send email")
		return err
	}

	logger.Log.WithFields(logrus.Fields{
		"to":     to,
		"subject": subject,
		"with_attachment": attachment != nil,
	}).Info("Email sent successfully")

	return nil
}