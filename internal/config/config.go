package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/killuox/koi/internal/shared"
	"gopkg.in/yaml.v2"
)

func Read(vars map[string]any) (cfg shared.Config, err error) {
	yamlFile, err := os.ReadFile("koi.config.yaml")

	// Regex to find {{variable}}
	re := regexp.MustCompile(`\{\{(\w+)\}\}`)

	// Replace all placeholders
	newYamlString := re.ReplaceAllStringFunc(string(yamlFile), func(match string) string {
		// Extract the key without {{}}
		key := strings.Trim(match, "{}")
		key = strings.TrimSpace(key)

		// Lookup the key in vars
		if val, ok := vars[key]; ok {
			return fmt.Sprintf("%v", val)
		}
		// If not found, keep original
		return match
	})

	// Convert back to byte
	yamlByte := []byte(newYamlString)

	if err != nil {
		return shared.Config{}, fmt.Errorf("error reading or missing koi.config.yaml file")
	}
	var config shared.Config
	err = yaml.Unmarshal(yamlByte, &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshalling config file")
	}

	return config, nil
}

func Validate(cfg shared.Config) error {

	return nil
}
