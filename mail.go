package golibs

import (
	"gopkg.in/gomail.v2"
)

func SendToMail(user,display_name, password, host, to, subject, body, mailtype string,port int) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", user, display_name)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/"+mailtype, body)
	d := gomail.NewDialer(host, port, user, password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
