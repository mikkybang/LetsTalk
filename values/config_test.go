package values

import (
	"testing"
)

func TestConfig(t *testing.T) {
	err := LoadConfiguration("../config.json")
	if err != nil {
		t.Error("could not load config", err)
	}

	if Config.Port == "" {
		t.Error("did not set host port")
	}
}
