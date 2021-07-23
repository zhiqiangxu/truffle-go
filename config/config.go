package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config ...
type Config struct {
	Node       string
	PrivateKey string
	ChainID    int
	Quick      bool
	Verbose    bool
}

// LoadConfig ...
func LoadConfig(confFile string) (config *Config, err error) {
	jsonBytes, err := ioutil.ReadFile(confFile)
	if err != nil {
		return
	}

	config = &Config{}
	err = json.Unmarshal(jsonBytes, config)
	return
}
