package values

import (
	"encoding/json"
	"os"
)

type configuration struct {
	DbHost                   string
	Port                     string
	CertPath                 string
	KeyPath                  string
	EnableClassSessionUpload bool // If EnableClassSessionUpload is set to true and no token is provided, files are saved to DB using GridFS.
	DropboxToken             string
}

// Config contains application environment variables.
var Config configuration

// LoadConfiguration loads all application environment variables.
func LoadConfiguration(configPath string) error {
	file, err := os.Open(configPath) // For read access.
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
