package checkprofile

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config object
type Config struct {
	Account *Account
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
		fmt.Println(err)
		return conf, err
	}
	if _, err = toml.DecodeFile(filename, &conf); err != nil {
		panic(err)
	}
	return conf, err
}
