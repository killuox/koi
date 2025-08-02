package commands

import (
	"fmt"
	"os"

	"github.com/killuox/koi/internal/api"
	"github.com/killuox/koi/internal/config"
	"gopkg.in/yaml.v2"
)

type State struct {
	config config.Config
}

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
		cfg, err := s.readConfig()
		if err != nil {
			return err
		}

		ep, ok := cfg.Endpoints[cmd.name]
		if !ok {
			return fmt.Errorf("Endpoint %s not found", cmd.name)
		}

		result, err := api.Call(ep, cfg)
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

func (s *State) readConfig() (cfg config.Config, err error) {
	yamlFile, err := os.ReadFile("koi.config.yaml")
	if err != nil {
		return config.Config{}, fmt.Errorf("Error reading or missing koi.config.yaml file")
	}
	var config config.Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, fmt.Errorf("Error unmarshalling config file")
	}

	return config, nil
}

// HANDLERS
func handlerHelp(s *State, cmd Command) error {

	return nil
}
