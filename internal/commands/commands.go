package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/killuox/koi/internal/api"
	"github.com/killuox/koi/internal/config"
	"gopkg.in/yaml.v2"
)

type State struct {
	cfg config.Config
}

type Command struct {
	name     string
	args     []string
	handler  commandHandler
	endpoint config.Endpoint
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

	cfg, err := state.readConfig()
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	state.cfg = cfg
	cName := os.Args[1]
	args := os.Args[2:]

	handler, ok := commands.handlers[cName]
	var endpoint config.Endpoint
	if !ok {

		ep, ok := cfg.Endpoints[cName]
		if !ok {
			fmt.Printf("no endpoints found for %s\n", cName)
			os.Exit(1)
		}

		endpoint = ep
	}

	cmd := Command{
		name:     cName,
		args:     args,
		handler:  handler,
		endpoint: endpoint,
	}

	err = commands.run(state, cmd)
	if err != nil {
		fmt.Printf("Error while running the command: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func (c *commands) run(s *State, cmd Command) error {
	if cmd.handler != nil {
		cmd.handler(s, cmd)
		return nil
	} else if cmd.endpoint.Path != "" {
		result, err := api.Call(cmd.endpoint, s.cfg)
		if err != nil {
			return fmt.Errorf("error while calling %s endpoint %s: %w", cmd.name, cmd.endpoint.Path, err)
		}

		var data interface{}
		unmarshalErr := json.Unmarshal(result, &data)
		if unmarshalErr == nil {
			prettyJSON, marshalErr := json.MarshalIndent(data, "", "  ")
			if marshalErr == nil {
				fmt.Print(string(prettyJSON))
			} else {
				fmt.Printf("Warning: Failed to pretty print JSON, printing raw result.\n%s\n", string(result))
			}
		} else {
			fmt.Printf("Result is not valid JSON, printing as plain text:\n%s\n", string(result))
		}
		return nil
	}

	return fmt.Errorf("Command %s not found", cmd.name)
}

func (c *commands) register(name string, f commandHandler) {
	c.handlers[name] = f
}

func (s *State) readConfig() (cfg config.Config, err error) {
	yamlFile, err := os.ReadFile("koi.config.yaml")
	if err != nil {
		return config.Config{}, fmt.Errorf("error reading or missing koi.config.yaml file")
	}
	var config config.Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshalling config file")
	}

	return config, nil
}

// HANDLERS
func handlerHelp(s *State, cmd Command) error {

	return nil
}
