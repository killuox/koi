package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-playground/validator/v10"
	"github.com/killuox/koi/internal/env"
	"gopkg.in/yaml.v2"
)

type Config struct {
	API       API                 `yaml:"api" validate:"required"`
	Endpoints map[string]Endpoint `yaml:"endpoints" validate:"required,dive"`
}

type API struct {
	BaseURL string            `yaml:"baseUrl" validate:"required,url"`
	Headers map[string]string `yaml:"headers"`
}

type SetVariableConfig struct {
	Body map[string]any `yaml:"body"`
}

type Endpoint struct {
	Method       string               `yaml:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Path         string               `yaml:"path" validate:"required"`
	Mode         string               `yaml:"mode" validate:"omitempty,oneof=env faker"`
	Parameters   map[string]Parameter `yaml:"parameters" validate:"dive"`
	Defaults     map[string]any       `yaml:"defaults"`
	SetVariables SetVariableConfig    `yaml:"set-variables"`
}

type Parameter struct {
	Type        string `yaml:"type" validate:"required,oneof=string int bool float"`
	Mode        string `yaml:"mode"`
	In          string `yaml:"in" validate:"omitempty,oneof=query path body"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Rules       Rules  `yaml:"rules"`
}

type Rules struct {
	// For strings
	MinLength int `yaml:"min_length" validate:"gte=0"`
	MaxLength int `yaml:"max_length" validate:"gte=0"`
	// Faker mode - Image
	Width  int `yaml:"width" validate:"gte=0"`
	Height int `yaml:"height" validate:"gte=0"`
	// Faker mode - For paragraph and sentence
	ParagraphCount int `yaml:"paragraph_count" validate:"gte=0"`
	SentenceCount  int `yaml:"sentence_count" validate:"gte=0"`
	WordCount      int `yaml:"word_count" validate:"gte=0"`
	// Faker mode - For numbers
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

// ENV
type EnvValueGetter interface {
	Get(key string, defaultVal any) (any, error)
}

type EnvStringParam struct{}
type EnvIntParam struct{}
type EnvBoolParam struct{}

var envParamTypeRegistry = map[string]EnvValueGetter{
	"string": EnvStringParam{},
	"int":    EnvIntParam{},
	"bool":   EnvBoolParam{},
}

// FAKER
type FakerValueGetter interface {
	Get(p Parameter) (any, error)
}

type FakerFullNameParam struct{}
type FakerFirstNameParam struct{}
type FakerLastNameParam struct{}
type FakerEmailParam struct{}
type FakerPasswordParam struct{}
type FakerCompanyParam struct{}
type FakerPhoneParam struct{}
type FakerLoremIpsumParam struct{}
type FakerNumberParam struct{}
type FakerImageParam struct{}
type FakerSentenceParam struct{}
type FakerParagraphParam struct{}

var fakerParamTypeRegistry = map[string]FakerValueGetter{
	"full_name":   FakerFullNameParam{},
	"first_name":  FakerFirstNameParam{},
	"last_name":   FakerLastNameParam{},
	"email":       FakerEmailParam{},
	"password":    FakerPasswordParam{},
	"company":     FakerCompanyParam{},
	"phone":       FakerPhoneParam{},
	"lorem_ipsum": FakerLoremIpsumParam{},
	"number":      FakerNumberParam{},
	"image":       FakerImageParam{},
	"sentence":    FakerSentenceParam{},
	"paragraph":   FakerParagraphParam{},
}

// Config
func (c *Config) Init(vars map[string]any) (err error) {
	yamlFile, err := os.ReadFile("koi.config.yaml")
	if err != nil {
		return fmt.Errorf("error reading koi.config.yaml file")
	}

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

	err = yaml.Unmarshal(yamlByte, c)
	if err != nil {
		return fmt.Errorf("error unmarshaling config file: %w", err)
	}

	return nil
}

func (c *Config) Validate(cfg Config) error {
	validate := validator.New()
	return validate.Struct(cfg)
}

func (c *Config) CreateValidatorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required but missing"
	case "url":
		return "must be a valid URL"
	case "oneof":
		return fmt.Sprintf("must be one of [%s]", e.Param())
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", e.Param())
	default:
		return fmt.Sprintf("failed validation rule '%s'", e.Tag())
	}
}

// Parameter
func (p Parameter) GetValue(flags map[string]any, key string, e Endpoint) (any, error) {
	// Get the default value
	defaultVal, hasDefaultValue := e.Defaults[key]

	// Check for flag value
	flagVal, ok := flags[key]
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

// ENV
func (EnvStringParam) Get(key string, defaultVal any) (any, error) {
	v, exists := env.GetString(key, "")
	if exists {
		return v, nil
	}

	return defaultVal, nil
}

func (EnvIntParam) Get(key string, defaultVal any) (any, error) {
	v, exists := env.GetInt(key, 0)
	if exists {
		return v, nil
	}

	return defaultVal, nil
}

func (EnvBoolParam) Get(key string, defaultVal any) (any, error) {
	v, exists := env.GetBool(key, false)
	if exists {
		return v, nil
	}

	return defaultVal, nil
}

func (p Parameter) GetEnvValue(key string, defaultVal any) (any, error) {
	getter, ok := envParamTypeRegistry[p.Type]
	if !ok {
		return nil, fmt.Errorf("unsupported parameter type: %s", p.Type)
	}

	return getter.Get(key, defaultVal)
}

// FAKER
func (FakerFullNameParam) Get(p Parameter) (any, error) {
	return gofakeit.Name(), nil
}
func (FakerFirstNameParam) Get(p Parameter) (any, error) {
	return gofakeit.FirstName(), nil
}
func (FakerLastNameParam) Get(p Parameter) (any, error) {
	return gofakeit.LastName(), nil
}
func (FakerEmailParam) Get(p Parameter) (any, error) {
	return gofakeit.Email(), nil
}
func (FakerPasswordParam) Get(p Parameter) (any, error) {
	return gofakeit.Password(true, true, true, true, false, 12), nil // TODO: make it dynamic via rules
}
func (FakerCompanyParam) Get(p Parameter) (any, error) {
	return gofakeit.Company(), nil
}
func (FakerPhoneParam) Get(p Parameter) (any, error) {
	return gofakeit.Phone(), nil
}
func (FakerLoremIpsumParam) Get(p Parameter) (any, error) {
	return gofakeit.LoremIpsumWord(), nil
}
func (FakerNumberParam) Get(p Parameter) (any, error) {
	return gofakeit.Number(p.Rules.Min, p.Rules.Max), nil
}
func (FakerImageParam) Get(p Parameter) (any, error) {
	return gofakeit.Image(p.Rules.Width, p.Rules.Height), nil
}
func (FakerSentenceParam) Get(p Parameter) (any, error) {
	return gofakeit.Sentence(p.Rules.WordCount), nil
}
func (FakerParagraphParam) Get(p Parameter) (any, error) {
	return gofakeit.Paragraph(p.Rules.ParagraphCount, p.Rules.SentenceCount, p.Rules.WordCount, "\n"), nil
}

func (p Parameter) GetFakerValue(key string) (any, error) {
	getter, ok := fakerParamTypeRegistry[key]
	if !ok {
		return nil, fmt.Errorf("%s is not a supported faker param", key)
	}

	return getter.Get(p)
}
