package config_test

import (
	"encoding/json"
	"testing"

	"github.com/driif/echo-go-starter/internal/server/config"
)

func TestPrintServiceEnv(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	_, err := json.MarshalIndent(config, "", "  ")

	if err != nil {
		t.Fatal(err)
	}
}
