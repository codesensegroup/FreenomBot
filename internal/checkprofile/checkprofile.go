package checkprofile

import (
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config object
type Config struct {
	Accounts []Account
	System   System
}

// System for basic auth
type System struct {
	Account    string
	Password   string
	CronTiming string
}

// Account struct
type Account struct {
	Username string
	Password string
}

// ReadConf will decode data
func ReadConf(filename string) (*Config, error) {
	var (
		conf *Config
		err  error
	)
	filename, err = filepath.Abs(filename)
	if err != nil {
		log.Fatal(err)
		return conf, err
	}
	if _, err = toml.DecodeFile(filename, &conf); err != nil {
		log.Fatal(err)
	}
	return conf, err
}
