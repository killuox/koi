package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/killuox/koi/internal/api"
	"github.com/killuox/koi/internal/config"
	"github.com/killuox/koi/internal/output"
	"gopkg.in/yaml.v2"
)

type State struct {
	cfg config.Config
}

type Command struct {
	name     string
	args     []string
	endpoint config.Endpoint
}

type commands struct {
}

func Init() {
	commands := commands{}

	state := &State{}

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

	ep, ok := cfg.Endpoints[cName]
	if !ok {
		fmt.Printf("no endpoints found for %s\n", cName)
		os.Exit(1)
	}

	cmd := Command{
		name:     cName,
		args:     args,
		endpoint: ep,
	}

	err = commands.run(state, cmd)
	if err != nil {
		fmt.Printf("Error while running the command: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func (c *commands) runWithLoader(
	f func() (api.Result, error),
) (api.Result, error) {
	done := make(chan struct{})
	var loaderProgram *tea.Program

	go func() {
		select {
		case <-done:
			// Loader should exit gracefully
			return
		case <-time.After(500 * time.Millisecond):
			// Initialize and run the loader if the operation takes too long
			loaderProgram = tea.NewProgram(output.InitLoader())
			if _, err := loaderProgram.Run(); err != nil {
				log.Printf("Error running loader: %v", err)
			}
		}
	}()

	// Execute the main function
	result, err := f()

	// Quit the loader if it was started
	if loaderProgram != nil {
		loaderProgram.Quit()
	}
	close(done) // Signal the goroutine to stop

	return result, err
}

// processAPIResult handles the API call result, including JSON formatting and output.
func (c *commands) processAPIResult(
	result api.Result,
) {
	var data interface{}
	unmarshalErr := json.Unmarshal(result.Body, &data)

	if unmarshalErr == nil {
		// Attempt to pretty print JSON
		prettyJSON, marshalErr := json.MarshalIndent(data, "", "  ")
		if marshalErr == nil {
			result.Body = prettyJSON
			output.ShowResponse(result)
		} else {
			fmt.Printf(
				"Warning: Failed to pretty print JSON, printing raw result.\n%s\n",
				string(result.Body),
			)
		}
	} else {
		// Not valid JSON, print as plain text
		fmt.Printf(
			"Result is not valid JSON, printing as plain text:\n%s\n",
			string(result.Body),
		)
	}
}

func (c *commands) run(s *State, cmd Command) error {
	callFunc := func() (api.Result, error) {
		startTime := time.Now()
		result, err := api.Call(cmd.endpoint, s.cfg)
		endTime := time.Now()

		if err == nil {
			result.Duration = endTime.Sub(startTime)
		}
		return result, err
	}

	result, err := c.runWithLoader(callFunc)
	if err != nil {
		return fmt.Errorf(
			"error while calling %s endpoint %s: %w",
			cmd.name,
			cmd.endpoint.Path,
			err,
		)
	}

	c.processAPIResult(result)
	return nil
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
