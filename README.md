# FreenomBot
The bot will auto-renew domain name
---
# How to use it

## Edit config.toml
Please type Freenom account(s)
``` toml
[System]
Account = "admin"
Password = "admin"
CronTiming = "23:30"

[[Accounts]]
Username = "xxx@xxx.com"
Password = "xxx"

[[Accounts]]
Username = "ooo@ooo.com"
Password = "ooo"
```

## Launch FreenomBot

On linux
``` sh
./FreenomBot_amd64
```
It will start http service on server, So you may check the status page of FrenomBot on http://localhost:8080

# People

Name: Teng-Wei, Hsieh

Mail: frank30941@gmail.com
