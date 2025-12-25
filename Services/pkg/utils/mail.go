package utils

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	smtpNet "net/smtp"
	"strings"

	"github.com/spf13/viper"
)

// SendMailModel represents email data
type SendMailModel struct {
	SenderName   string
	SenderMail   string
	ReceiverName string
	ReceiverMail string
	Subject      string
	Body         string
	AltBody      string
	Username     string
	Password     string
}

// SendNormalMail sends a basic email
func SendNormalMail(m SendMailModel) (bool, error) {
	if m.ReceiverMail == "" {
		return false, nil
	}
	if m.Body == "" {
		m.Body = m.AltBody
	}

	smtpHost := viper.GetString("app.smtp.host_test")
	smtpPort := viper.GetInt("app.smtp.port_test")
	m.Username = viper.GetString("app.smtp.username_test")
	m.Password = viper.GetString("app.smtp.password_test")
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	from := mail.Address{Name: m.SenderName, Address: m.SenderMail}
	to := mail.Address{Name: m.ReceiverName, Address: m.ReceiverMail}

	// Build message
	header := []string{
		"From: " + from.String(),
		"To: " + to.String(),
		"Subject: " + m.Subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"utf-8\"",
		"Content-Transfer-Encoding: base64",
	}
	var msg strings.Builder
	for _, h := range header {
		msg.WriteString(h + "\r\n")
	}
	msg.WriteString("\r\n")
	msg.WriteString(base64.StdEncoding.EncodeToString([]byte(m.Body)))

	// AUTH
	auth := smtpNet.PlainAuth("", m.Username, m.Password, smtpHost)

	// Send
	if err := smtpNet.SendMail(addr, auth, from.Address, []string{to.Address}, []byte(msg.String())); err != nil {
		return false, err
	}
	return true, nil
}
