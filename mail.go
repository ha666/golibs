package golibs

import (
	"net/smtp"
	"strings"
)

func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	auth := smtp.PlainAuth("", user, password, strings.Split(host, ":")[0])
	msg := []byte("To: " + to + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\n" + "Content-Type: text/" + mailtype + "; charset=UTF-8" + "\r\n\r\n" + body)
	sendto := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, sendto, msg)
	return err
}
