package i18ntpl

import (
	"log"

	"github.com/codesensegroup/FreenomBot/internal/config"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

//Init create bundle
func Init(cfg *config.Config) {
	log.Fatalln("Init i18n packge.")
	bundle = i18n.NewBundle(language.MustParse(cfg.System.Lang))
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFile("i18n/mail.notion." + cfg.System.Lang + ".toml")
	log.Fatalln("Finish i18n packge.")
}
