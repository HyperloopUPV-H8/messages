package main

import (
	"os"
	"strings"

	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Excel excel_adapter.ExcelAdapterConfig
}

func getConfig(path string) (Config, error) {
	configFile, err := os.ReadFile(path)

	if err != nil {
		return Config{}, err
	}

	reader := strings.NewReader(string(configFile))

	var config Config

	err = toml.NewDecoder(reader).Decode(&config)
	return config, err
}
