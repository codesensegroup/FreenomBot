package mailer

import (
	gomail "gopkg.in/gomail.v2"
)

func SendMail(account, password, to string, msg *string) {
	m := gomail.NewMessage()
	m.SetHeader("From", account)
	m.SetHeader("To", to)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "FreenomBot")
	m.SetBody("text/html", *msg)

	d := gomail.NewDialer("smtp.example.com", 587, account, password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
