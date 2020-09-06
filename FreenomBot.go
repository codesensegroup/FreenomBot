package main

import (
	"log"

	checkprofile "github.com/codesensegroup/FreenomBot/internal/checkprofile"
	freenom "github.com/codesensegroup/FreenomBot/internal/freenom"
	httpservice "github.com/codesensegroup/FreenomBot/server/httpservice"

	"github.com/robfig/cron"
)

func runFreenom(run func()) {
	c := cron.New()
	c.AddFunc("* * 23 */1 * *", run)
	c.Start()
	select {}
}

func main() {
	log.Println("init")
	config, _ := checkprofile.ReadConf("./config.toml")
	f := freenom.GetInstance()

	go func() {
		for i, a := range config.Accounts {
			f.InputAccount(i, a.Username, a.Password)
			f.Login(i).RenewDomains(i)
			for _, d := range f.Users[i].Domains {
				log.Println("log: ", d)
			}
			runFreenom(func() {
				f.Login(i).RenewDomains(i)
				for _, d := range f.Users[i].Domains {
					log.Println("log: ", d)
				}
			})
		}
	}()

	httpservice.Run(f)
}
