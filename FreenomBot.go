package main

import (
	"log"

	"github.com/codesensegroup/FreenomBot/internal/checkprofile"
	"github.com/codesensegroup/FreenomBot/internal/freenom"
	"github.com/codesensegroup/FreenomBot/internal/scheduler"
	"github.com/codesensegroup/FreenomBot/server/httpservice"
)

func task(f *freenom.Freenom, acs int) {
	var i int
	for i = 0; i < acs; i++ {
		f.Login(i).RenewDomains(i)
		for _, d := range f.Users[i].Domains {
			log.Println("log: ", d)
		}
	}
}

func main() {
	log.Println("Init")
	config, _ := checkprofile.ReadConf("./config.toml")
	f := freenom.GetInstance()
	f.InputAccount(config)
	task(f, len(config.Accounts))
	go scheduler.Run(func() {
		task(f, len(config.Accounts))
	}, config.System.CronTiming)

	httpservice.Run(f, config)
}
