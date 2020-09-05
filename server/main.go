package main

import (
	"log"

	checkprofile "github.com/codesensegroup/FreenomBot/internal/checkprofile"
	freenom "github.com/codesensegroup/FreenomBot/internal/freenom"

	"github.com/robfig/cron"
)

func main() {
	log.Println("init")
	config, _ := checkprofile.ReadConf("./config.toml")
	c := cron.New()
	f := freenom.GetInstance()
	c.AddFunc("*/5 * * * * *", func() {
		f.Login(config.Account.Username, config.Account.Password).RenewDomains()
		for _, d := range f.Domains {
			log.Println("log: ", d)
		}
	})
	c.Start()
	select {}
}
