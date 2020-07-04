package values

import (
	"encoding/json"
	"log"
	"os"
)

type configuration struct {
	DbHost string
	Port   string
}

var Config configuration

func LoadConfiguration() error {
	file, err := os.Open("./config.json") // For read access.
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		return err
	}
	return nil
}
