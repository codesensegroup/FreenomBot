package mailer

import (
	"bytes"
	"html/template"
	"log"

	"github.com/codesensegroup/FreenomBot/internal/freenom"
	gomail "gopkg.in/gomail.v2"
)

func SendMail(f *freenom.Freenom) {

	m := gomail.NewMessage()
	m.SetHeader("From", f.ConfigData.Mailer.Account)
	m.SetHeader("To", f.ConfigData.Mailer.To)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "FreenomBot")
	m.SetBody("text/html", getMailText(f))

	d := gomail.NewDialer("smtp.example.com", 587, f.ConfigData.Mailer.Account, f.ConfigData.Mailer.Password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func getMailText(f *freenom.Freenom) (text string) {
	tmplate := "notice"
	var tmp bytes.Buffer
	t, err := template.ParseFiles("./resources/mail/" + tmplate + ".html")
	if err != nil {
		log.Fatalln("error ./resources/mail/" + tmplate + ".html")
	}
	err = t.Execute(&tmp, f.ConfigData)
	text = string(tmp.Bytes()[:])
	return
}
