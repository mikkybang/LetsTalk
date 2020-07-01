package model

import (
	"encoding/json"
	"log"
	"os"
)

var Config Configuration

func LoadConfiguration() {
	file, err := os.Open("config.json") // For read access.
	if err != nil {
		log.Fatal("Error loading the config file")
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("can't decode config JSON: ", err)
	}
}
