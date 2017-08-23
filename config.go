package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

const defaultBind = ":9090"

type Config struct {
	Bind            string
	OFBSymmetricKey string
	ApiKey          string
}

func LoadConfig(path string) (*Config, error) {
	rawBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	err = json.Unmarshal(rawBytes, &config)
	if err != nil {
		return nil, err
	}

	if len(config.OFBSymmetricKey) == 0 {
		return nil, errors.New("Config is missing a OFBSymmetricKey")
	}

	if len(config.Bind) == 0 {
		config.Bind = defaultBind
	}

	return config, nil
}
