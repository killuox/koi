package config

import (
	"fmt"
	"os"

	"github.com/killuox/koi/internal/shared"
	"gopkg.in/yaml.v2"
)

func Read() (cfg shared.Config, err error) {
	yamlFile, err := os.ReadFile("koi.config.yaml")
	if err != nil {
		return shared.Config{}, fmt.Errorf("error reading or missing koi.config.yaml file")
	}
	var config shared.Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshalling config file")
	}

	return config, nil
}
