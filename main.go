package main

import (
	"log"

	"FreenomBot/common/freenom"
	"FreenomBot/common/line"
	"FreenomBot/common/message"
	"FreenomBot/common/scheduler"
	"FreenomBot/config"
	"FreenomBot/server/httpservice"
)

func task(f *freenom.Freenom, pageData *freenom.PageData) {
	var err error
	var isRenew bool
	for i := range pageData.Users {
		if err = f.Login(&pageData.Users[i]); err != nil {
			log.Fatalln(err)
			continue
		}

		if err = f.GetFreenomInfo(&pageData.Users[i]); err != nil {
			log.Fatalln(err)
			continue
		}

		if err = f.RenewDomains(&pageData.Users[i]); err != nil {
			log.Fatalln(err)
			continue
		}

		for _, d := range pageData.Users[i].Domains {
			log.Println("token: ", pageData.Users[i].Token)
			log.Println("domain: ", d)
			if d.RenewState != freenom.RenewNo {
				isRenew = true
			}
		}
	}

	if err != nil {
		log.Fatalln(err)
		return
	}

	msg, err := message.GenMessage(pageData)
	if err != nil {
		log.Fatalln(err)
		return
	}

	if configData.Line.Enable {
		if configData.Line.Daily {
			line.Send(&configData.Line.Token, msg)
		} else if isRenew {
			line.Send(&configData.Line.Token, msg)
		}
	}

}

var configData *config.Config

func main() {
	log.Println("Init")
	configData = config.GetData()
	//i18nTpl.Init(configData)
	if configData.Line.Enable {
		line.Init()
	}

	pageData := &freenom.PageData{
		Users: make([]freenom.User, len(configData.Accounts)),
	}

	for i := range configData.Accounts {
		pageData.Users[i].UserName = configData.Accounts[i].Username
		pageData.Users[i].Password = configData.Accounts[i].Password
	}

	f := freenom.GetInstance()
	task(f, pageData)
	go scheduler.Run(func() {
		task(f, pageData)
	}, configData.System.CronTiming)

	httpservice.Run(pageData, configData)
}
