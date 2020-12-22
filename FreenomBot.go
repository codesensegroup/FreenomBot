package main

import (
	"fmt"
	"log"

	"github.com/codesensegroup/FreenomBot/internal/config"
	"github.com/codesensegroup/FreenomBot/internal/freenom"
	"github.com/codesensegroup/FreenomBot/internal/line"
	"github.com/codesensegroup/FreenomBot/internal/scheduler"
	"github.com/codesensegroup/FreenomBot/server/httpservice"
)

func task(f *freenom.Freenom, acs int) {
	var i int
	for i = 0; i < acs; i++ {
		f.Login(i).RenewDomains(i)
		for _, d := range f.ConfigData.Accounts[i].Domains {
			log.Println("log: ", d)
			line.Send(fmt.Sprintf("%#v", d))
		}
	}
}

func main() {
	log.Println("Init")
	configData := config.GetData()
	//i18nTpl.Init(configData)
	if configData.Line.Enable {
		line.Init(&configData.Line.Token)
	}
	f := freenom.GetInstance()
	task(f, len(configData.Accounts))
	go scheduler.Run(func() {
		task(f, len(configData.Accounts))
	}, configData.System.CronTiming)

	httpservice.Run(f, configData)
}
