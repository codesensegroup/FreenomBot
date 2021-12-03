package config

import (
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
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
	Lang       string
}

// Mailer for send email at tom file
type Mailer struct {
	Enable   bool
	Daily    bool
	Account  string
	Password string
	To       string
}

// Line for send message at tom file
type Line struct {
	Enable bool
	Daily  bool
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
	if err := readConf("./config.toml"); err != nil {
		log.Fatalln(err)
	}
}
