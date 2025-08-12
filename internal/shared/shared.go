package shared

import "fmt"

type State struct {
	Cfg   Config
	Flags map[string]any
}

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

func (p Parameter) GetValue(s *State, key string, e Endpoint) (any, error) {
	// Check for flag value
	flagVal, ok := s.Flags[key]
	if ok {
		return flagVal, nil
	}
	// Get the default value
	defaultVal, ok := e.Defaults[key]
	if ok {
		return defaultVal, nil
	}
	return nil, fmt.Errorf("no value provided for parameter: %s", key)
}
