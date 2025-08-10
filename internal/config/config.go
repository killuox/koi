package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	API       API                 `yaml:"api"`
	Endpoints map[string]Endpoint `yaml:"endpoints"`
}

type API struct {
	BaseURL string `yaml:"baseUrl"`
	Version string `yaml:"version"`
	Auth    Auth   `yaml:"auth"`
}

type Auth struct {
	Type string `yaml:"type"`
}

type Endpoint struct {
	Type       string               `yaml:"type"`
	Method     string               `yaml:"method"`
	Mode       string               `yaml:"mode"`
	Path       string               `yaml:"path"`
	Parameters map[string]Parameter `yaml:"parameters"`
	Defaults   map[string]any       `yaml:"defaults"`
}

type Parameter struct {
	Type        string     `yaml:"type"`
	Mode        string     `yaml:"mode"`
	In          string     `yaml:"in"`
	Description string     `yaml:"description"`
	Required    bool       `yaml:"required"`
	Validation  Validation `yaml:"validation"`
}

type Validation struct {
	MinLength int `yaml:"minLength"`
	MaxLength int `yaml:"maxLength"`
}

func (p Parameter) GetValue(key string, e Endpoint) (any, error) {
	// Check for flag value
	flagVal, err := p.GetFlagValue(key)
	if err == nil {
		return flagVal, nil
	}
	// Get the default value
	defaultVal, ok := e.Defaults[key]
	if ok {
		return defaultVal, nil
	}
	return nil, fmt.Errorf("no value found for parameter: %s", key)
}

func (p Parameter) GetFlagValue(key string) (any, error) {
	epName := os.Args[1]
	cmd := flag.NewFlagSet(epName, flag.ExitOnError)
	keyFLag := cmd.Lookup(key)
	if keyFLag == nil {
		return nil, fmt.Errorf("no flag value found for %s", key)
	}

	var value interface{}
	switch p.Type {
	case "string":
		value = p.getStringValue(key, cmd)
	case "bool":
		value = p.getBoolValue(key, cmd)
	case "int":
		value = p.getIntValue(key, cmd)
	default:
		value = nil
	}

	cmd.Parse(os.Args[2:])

	switch v := value.(type) {
	case string:
		if v != "" {
			return v, nil
		}
	case bool:
		// Accept both true and false
		return v, nil
	case int:
		if v != 0 {
			return v, nil
		}
	default:
		return nil, fmt.Errorf("unsupported flag type: %T", v)
	}

	return nil, fmt.Errorf("no flag value found for %s", key)
}

func (p Parameter) getStringValue(key string, cmd *flag.FlagSet) string {
	return *cmd.String(key, "", "")
}

func (p Parameter) getBoolValue(key string, cmd *flag.FlagSet) bool {
	return *cmd.Bool(key, false, "")
}

func (p Parameter) getIntValue(key string, cmd *flag.FlagSet) int {
	return *cmd.Int(key, 0, "")
}
