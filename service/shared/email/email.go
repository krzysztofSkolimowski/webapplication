package email

import (
	"net/smtp"
	"fmt"
	"encoding/base64"
)

var (
	e SMTPInfo
)

type SMTPInfo struct {
	Username string
	Pass     string
	Host     string
	Port     int
	From     string
}

func Configure(c SMTPInfo) {
	e = c
}

func ReadConfig() SMTPInfo {
	return e
}

func SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", e.Username, e.Pass, e.Host)

	header := make(map[string]string)
	header["From"] = e.From
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = `text/plain; charset="utf-8"`
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", e.Host, e.Port),
		auth,
		e.From,
		[]string{to},
		[]byte(message),
	)

	return err
}
