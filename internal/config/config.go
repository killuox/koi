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

func (e Endpoint) GetValue() string {
	cmd := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	slug := cmd.String("slug", "", "The slug of the Pokemon to retrieve.")
	cmd.Parse(os.Args[2:])

	if *slug == "" {
		fmt.Println("Error: --slug is required for get-pokemon command.")
		cmd.Usage()
		os.Exit(1)
	}

	return *slug
}
