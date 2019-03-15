package main

import (
	"os"
	"errors"
	"encoding/json"
)

var (
	CURRENT_PATH, _ = os.Getwd()
)

type Config struct {
	Port int            `json:"port,omitempty"`
	MQ   *MQConfig      `json:"mq,omitempty"`
}

type MQConfig struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	InputQueue     string `json:"input_queue"`
	OutputExchange string `json:"output_exchange"`
}

func loadConfig(configFileName string) (*Config, error) {
	// Config object to set defaults
	config := &Config{
		Port: 8765,
		MQ:             &MQConfig{},
	}

	configFile := CURRENT_PATH + string(os.PathSeparator) + configFileName

	file, err := os.Open(configFile)
	if err != nil {
		return config, errors.New("Couldn't open configfile %s with err " + configFile + err.Error())
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, errors.New("Couldn't parse configfile %s with err " + configFile + err.Error())
	}

	return config, nil
}

