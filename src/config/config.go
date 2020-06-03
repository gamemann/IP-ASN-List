package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Token      string `json:"token"`
	Port       int    `json:"port"`
	UpdateTime int    `json:"updatetime"`
}

func ReadConfig(cfg *Config, filename string) bool {
	// Open config file.
	file, err := os.Open(filename)

	// Check for errors.
	if err != nil {
		fmt.Println("Error opening config file.")

		return false
	}

	// Defer file close.
	defer file.Close()

	// Create stat.
	stat, _ := file.Stat()

	// Make data.
	data := make([]byte, stat.Size())

	// Read data.
	_, err = file.Read(data)

	// Check for errors.
	if err != nil {
		fmt.Println("Error reading config file.")

		return false
	}

	// Parse JSON data.
	err = json.Unmarshal([]byte(data), cfg)

	// Check for errors.
	if err != nil {
		fmt.Println("Error parsing JSON Data.")

		return false
	}

	return true
}
