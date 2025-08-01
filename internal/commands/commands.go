package commands

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type State struct {
	config Config
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
	Method     string                 `yaml:"method"`
	Mode       string                 `yaml:"mode"` // This field is optional and only present in some endpoints.
	Path       string                 `yaml:"path"`
	Parameters map[string]Parameter   `yaml:"parameters"`
	Defaults   map[string]interface{} `yaml:"defaults"`
}

type Parameter struct {
	Type        string     `yaml:"type"`
	Mode        string     `yaml:"mode"` // This field is optional and only present in some parameters.
	In          string     `yaml:"in"`
	Description string     `yaml:"description"`
	Required    bool       `yaml:"required"`
	Validation  Validation `yaml:"validation"` // This field is optional.
}

type Validation struct {
	MinLength int `yaml:"minLength"`
	MaxLength int `yaml:"maxLength"`
}

type Result struct{}

type Command struct {
	name    string
	args    []string
	handler commandHandler
}

type commandHandler func(s *State, cmd Command) error

type commands struct {
	handlers map[string]commandHandler
}

func Init() {
	commands := commands{
		handlers: make(map[string]commandHandler),
	}

	state := &State{}

	// register commands here
	commands.register("help", handlerHelp)

	if len(os.Args) < 2 {
		fmt.Print("Not enough arguments provided.\n")
		os.Exit(1)
	}

	cName := os.Args[1]
	args := os.Args[2:]

	handler, ok := commands.handlers[cName]
	if !ok {
		fmt.Printf("Command name '%s' not found\n", cName)
		os.Exit(1)
	}

	cmd := Command{
		name:    cName,
		args:    args,
		handler: handler,
	}

	err := commands.run(state, cmd)
	if err != nil {
		fmt.Printf("Error while running the command: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func (c *commands) run(s *State, cmd Command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		ep, err := findEndpoint(s, cmd)
		if err != nil {
			return err
		}

		result, err := ep.call(s)
		if err != nil {
			return err
		}

		//TODO: Display the result beautifully in the terminal
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f commandHandler) {
	c.handlers[name] = f
}

func readConfig() (cfg Config, err error) {
	yamlFile, err := os.ReadFile("koi.config.yaml")
	if err != nil {
		return Config{}, fmt.Errorf("Error reading or missing koi.config.yaml file")
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Error unmarshalling config file")
	}

	return config, nil
}

func findEndpoint(s *State, cmd Command) (e Endpoint, err error) {
	cfg, err := readConfig()
	if err != nil {
		return Endpoint{}, err
	}

	// set cfg to the state
	s.config = cfg

	ep, ok := cfg.Endpoints[cmd.name]
	if !ok {
		return Endpoint{}, fmt.Errorf("Endpoint %s not found", cmd.name)
	}

	return ep, nil
}

// HANDLERS
func handlerHelp(s *State, cmd Command) error {

	return nil
}

func (e *Endpoint) call(s *State) (r Result, err error) {
	//TODO make a http call
	return Result{}, nil
}
