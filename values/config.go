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

// Config contains application environment variables.
var Config configuration

func init() {
	err := LoadConfiguration()
	if err != nil {
		log.Fatalln("could not load config", err)
	}
}

// LoadConfiguration loads all application environment variables.
func LoadConfiguration() error {
	file, err := os.Open("../config.json") // For read access.
	if err != nil {
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
