package values

import (
	"encoding/json"
	"log"
	"os"
)

var Config Configuration

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
