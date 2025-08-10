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
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
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
	if len(os.Args) < 2 {
		return nil, fmt.Errorf("no command provided")
	}

	epName := os.Args[1]
	cmd := flag.NewFlagSet(epName, flag.ExitOnError)

	var value interface{}
	switch p.Type {
	case "string":
		value = cmd.String(key, "", "string flag")
	case "bool":
		value = cmd.Bool(key, false, "bool flag")
	case "int":
		value = cmd.Int(key, 0, "int flag")
	default:
		return nil, fmt.Errorf("unsupported flag type: %s", p.Type)
	}

	// Parse the flags from the remaining arguments
	if err := cmd.Parse(os.Args[2:]); err != nil {
		return nil, err
	}

	// Return the dereferenced value
	switch v := value.(type) {
	case *string:
		if *v != "" {
			return *v, nil
		}
	case *bool:
		return *v, nil
	case *int:
		if *v != 0 {
			return *v, nil
		}
	}

	return nil, fmt.Errorf("no flag value found for %s", key)
}
