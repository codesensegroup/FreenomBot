package main

import (
	"log"

	checkprofile "github.com/codesensegroup/FreenomBot/internal/checkprofile"
	freenom "github.com/codesensegroup/FreenomBot/internal/freenom"
	httpservice "github.com/codesensegroup/FreenomBot/server/httpservice"

	//"github.com/robfig/cron"
	"github.com/jasonlvhit/gocron"
)

var cronJobs = make(map[int]*gocron.Scheduler)

func runFreenom(run func(id int), id int, timing string) {
	cronJobs[id] = gocron.NewScheduler()
	cronJobs[id].Every(1).Day().At(timing).Do(run, id)
	<-cronJobs[id].Start()
}

func main() {
	log.Println("init")
	config, _ := checkprofile.ReadConf("./config.toml")
	f := freenom.GetInstance()

	for i, a := range config.Accounts {
		f.InputAccount(i, a.Username, a.Password)
		f.Login(i).RenewDomains(i)
		for _, d := range f.Users[i].Domains {
			log.Println("log: ", d)
		}
		//Use goroutine
		go runFreenom(func(id int) {
			f.Login(id).RenewDomains(id)
			for _, d := range f.Users[id].Domains {
				log.Println("log: ", d)
			}
		}, i, config.System.CronTiming)
	}

	httpservice.Run(f, config)
}
