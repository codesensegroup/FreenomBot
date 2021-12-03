package message

import (
	"FreenomBot/common/freenom"
	"bytes"
	"text/template"
)

func GenMessage(f *freenom.PageData) (text string, err error) {
	tmplate := "notice"
	var tmp bytes.Buffer
	var t *template.Template
	if t, err = template.ParseFiles("./resources/mail/" + tmplate + ".tpl"); err != nil {
		return
	} else if err = t.Execute(&tmp, f); err != nil {
		return
	}

	text = string(tmp.Bytes()[:])

	return
}
