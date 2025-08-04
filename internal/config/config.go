package config

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
	Method     string                 `yaml:"method"`
	Mode       string                 `yaml:"mode"`
	Path       string                 `yaml:"path"`
	Parameters map[string]Parameter   `yaml:"parameters"`
	Defaults   map[string]interface{} `yaml:"defaults"`
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
