package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Connection struct {
		LocalAddress  string `json:"local_address"`
		RemoteAddress string `json:"remote_address"`
	} `json:"connection json:remote_address"`
}

func readConfig() Config {
	c := Config{}
	configPath := "config.json"

	c.Connection.LocalAddress = "0.0.0.0:19132"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		file, createErr := os.Create(configPath)
		if createErr != nil {
			fmt.Errorf("ERROR #1 CONFIG, exit code 1 | failed to create %s: %v", configPath, createErr)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if encodeErr := encoder.Encode(c); encodeErr != nil {
			fmt.Errorf("ERROR #2 CONFIG, exit code 1 | failed to write default configuration: %v", encodeErr)
		}
		os.Exit(0)
	}

	file, err := os.Open(configPath)
	if err != nil {
		fmt.Errorf("ERROR #3 CONFIG, exit code 1 | access denied to %s: %v", configPath, err)
	}
	defer file.Close()

	if decodeErr := json.NewDecoder(file).Decode(&c); decodeErr != nil {
		fmt.Errorf("ERROR #4 CONFIG, exit code 1 | failed to decode %s. check JSON syntax: %v", configPath, decodeErr)
	}

	return c
}
