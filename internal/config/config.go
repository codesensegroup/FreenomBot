package config

import (
	"log"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/BurntSushi/toml"
)

const (
	Version                 = "v0.0.5"
	Timeout                 = 34
	DelayRegistionDelayTime = 14
	BaseURL                 = "https://my.freenom.com"
	RefererURL              = "https://my.freenom.com/clientarea.php"
	LoginURL                = "https://my.freenom.com/dologin.php"
	DomainStatusURL         = "https://my.freenom.com/domains.php?a=renewals"
	RenewDomainURL          = "https://my.freenom.com/domains.php?submitrenewals=true"
	AuthKey                 = "WHMCSZH5eHTGhfvzP"
)

var (
	TokenREGEX       = regexp.MustCompile(`name="token"\svalue="(?P<token>[^"]+)"`)
	DomainInfoREGEX  = regexp.MustCompile(`<tr><td>(?P<domain>[^<]+)<\/td><td>[^<]+<\/td><td>[^<]+<span class="[^"]+">(?P<days>\d+)[^&]+&domain=(?P<id>\d+)"`)
	LoginStatusREGEX = regexp.MustCompile(`<li.*?Logout.*?<\/li>`)
	CheckRenew       = regexp.MustCompile(`(?i)Order Confirmation`)
)

const (
	RenewNo  int = 0
	RenewYes int = 1
	RenewErr int = 3
)

/** Tom file **/
// Config object at tom file
type Config struct {
	Accounts []Account
	System   System
	Mailer   Mailer
	Line     Line
}

// System for basic auth at tom file
type System struct {
	Account    string
	Password   string
	CronTiming string
	Daily      bool
	Lang       string
}

// Mailer for send email at tom file
type Mailer struct {
	Enable   bool
	Account  string
	Password string
	To       string
}

// Line for send message at tom file
type Line struct {
	Enable bool
	Token  string
}

// Account struct at tom file
type Account struct {
	Username string
	Password string
	Domains  []Domain
}

// Domain data
type Domain struct {
	DomainName string
	Days       int
	ID         string
	RenewState int
	CheckTimes int
}

var configData *Config
var once sync.Once

func readConf(filename string) error {
	var (
		err error
	)
	filename, err = filepath.Abs(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if _, err = toml.DecodeFile(filename, &configData); err != nil {
		log.Fatal(err)
	}
	return err
}

//CreateDomains is a creation Domains function
func CreateDomains(uid int, size int) {
	length := len(configData.Accounts[uid].Domains)
	if length != size {
		configData.Accounts[uid].Domains = make([]Domain, size)
	}
}

//GetData get Config data
func GetData() *Config {
	return configData
}

func init() {
	once.Do(func() {
		readConf("./config.toml")
	})
}
