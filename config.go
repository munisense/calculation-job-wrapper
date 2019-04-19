package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var (
	currentPath, _ = os.Getwd()
)

type Config struct {
	Port int       `json:"port,omitempty"`
	MQ   *MQConfig `json:"mq,omitempty"`
}

type MQConfig struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	InputQueue     string `json:"input_queue"`
	OutputExchange string `json:"output_exchange"`
}

func loadConfig() (*Config, error) {
	// Config object to set defaults
	config := &Config{
		Port: StringToIntUnsafe(os.Getenv("PORT"), 8080),
		MQ: &MQConfig{
			Host:           os.Getenv("MQ_HOST"),
			Port:           StringToIntUnsafe(os.Getenv("MQ_PORT"), 0),
			InputQueue:     os.Getenv("MQ_INPUT_QUEUE"),
			OutputExchange: os.Getenv("MQ_OUTPUT_EXCHANGE"),
			Username:       os.Getenv("MQ_USERNAME"),
			Password:       os.Getenv("MQ_PASSWORD"),
		},
	}

	configPath := currentPath + string(os.PathSeparator) + "config" + string(os.PathSeparator) + "config.json"
	file, err := os.Open(configPath)
	if err != nil {
		return config, fmt.Errorf("couldn't open configfile %s with err: %s", configPath, err.Error())
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, fmt.Errorf("couldn't decode configfile %s with err: %s", configPath, err.Error())
	}

	return config, nil
}

func StringToIntUnsafe(input string, defaultVal int) int {
	if input == "" {
		return defaultVal
	}

	output, err := strconv.Atoi(input)
	if err != nil {
		panic(err)
	}

	return output
}
