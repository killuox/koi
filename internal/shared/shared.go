package shared

import (
	"fmt"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/killuox/koi/internal/env"
)

var validFakerKeys = []string{"full_name", "first_name", "last_name", "email", "password", "company", "phone", "lorem_ipsum", "number", "image", "sentence", "paragraph"}

type State struct {
	Cfg   Config
	Flags map[string]any
}

type Config struct {
	API       API                 `yaml:"api"`
	Endpoints map[string]Endpoint `yaml:"endpoints"`
}

type API struct {
	BaseURL string            `yaml:"baseUrl"`
	Headers map[string]string `yaml:"headers"`
}

type Headers struct {
	Authorization string `yaml:"Authorization"`
	ContentType   string `yaml:"ContentType"`
}

// First key is the location
// the string key of the map is name you want to register the variable has
// the value is the key you want to get the value from
type SetVariableConfig struct {
	Body map[string]any `yaml:"body"`
}

type Endpoint struct {
	Type         string               `yaml:"type"`
	Method       string               `yaml:"method"`
	Mode         string               `yaml:"mode"`
	Path         string               `yaml:"path"`
	Parameters   map[string]Parameter `yaml:"parameters"`
	Defaults     map[string]any       `yaml:"defaults"`
	SetVariables SetVariableConfig    `yaml:"set-variables"`
}

type Parameter struct {
	Type        string `yaml:"type"`
	Mode        string `yaml:"mode"`
	In          string `yaml:"in"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Rules       Rules  `yaml:"rules"`
}

type Rules struct {
	// For strings
	MinLength int `yaml:"min_length"`
	MaxLength int `yaml:"max_length"`
	// For Image
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
	// For paragraph and sentence
	ParagraphCount int `yaml:"paragraph_count"`
	SentenceCount  int `yaml:"sentence_count"`
	WordCount      int `yaml:"word_count"`
	// For numbers
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

func (p Parameter) GetValue(s *State, key string, e Endpoint) (any, error) {
	// Get the default value
	defaultVal, hasDefaultValue := e.Defaults[key]

	// Check for flag value
	flagVal, ok := s.Flags[key]
	if ok {
		return flagVal, nil
	}

	modeParts := strings.Split(p.Mode, ":")
	if len(modeParts) == 2 {
		modeType := modeParts[0]
		modeValue := modeParts[1]

		if modeType == "env" {
			v, err := p.GetEnvValue(modeValue, defaultVal)
			if err == nil && v != "" {
				return v, nil
			}
		}

		if modeType == "faker" {
			v, err := p.GetFakerValue(modeValue)
			if err == nil {
				return v, nil
			}
		}
	}

	if hasDefaultValue {
		return defaultVal, nil
	}

	return nil, fmt.Errorf("no value provided for parameter: %s", key)
}

func (p Parameter) GetEnvValue(key string, defaultVal any) (any, error) {
	switch p.Type {
	case "string":
		v, exists := env.GetString(key, "")
		if exists {
			return v, nil
		}
	case "int":
		v, exists := env.GetInt(key, 0)
		if exists {
			return v, nil
		}
	case "bool":
		v, exists := env.GetBool(key, false)
		if exists {
			return v, nil
		}
	default:
		return nil, fmt.Errorf("wrong parameter type for: %s", key)
	}
	return defaultVal, nil
}

func (p Parameter) GetFakerValue(key string) (any, error) {
	var isValid bool
	for _, fk := range validFakerKeys {
		if fk == key {
			isValid = true
		}
	}
	if !isValid {
		return nil, fmt.Errorf("invalid key for faker mode make sure to provide a valid one")
	}

	switch key {
	case "first_name":
		return gofakeit.FirstName(), nil
	case "last_name":
		return gofakeit.LastName(), nil
	case "full_name":
		return gofakeit.Name(), nil
	case "phone":
		return gofakeit.Phone(), nil
	case "email":
		return gofakeit.Email(), nil
	case "password":
		return gofakeit.Password(true, true, true, true, false, 12), nil // TODO: make it dynamic via rules
	case "company":
		return gofakeit.Company(), nil
	case "image":
		return gofakeit.Image(p.Rules.Width, p.Rules.Height), nil
	case "number":
		return gofakeit.Number(p.Rules.Min, p.Rules.Max), nil
	case "sentence":
		return gofakeit.Sentence(p.Rules.WordCount), nil
	case "paragraph":
		return gofakeit.Paragraph(p.Rules.ParagraphCount, p.Rules.SentenceCount, p.Rules.WordCount, "\n"), nil
	default:
		return nil, fmt.Errorf("key %s does not exist in mode faker", key)
	}
}
