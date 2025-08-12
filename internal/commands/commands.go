package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/killuox/koi/internal/api"
	"github.com/killuox/koi/internal/config"
	"github.com/killuox/koi/internal/output"
	"github.com/killuox/koi/internal/shared"
)

type Command struct {
	name     string
	args     []string
	endpoint shared.Endpoint
}

type commands struct {
}

func Init() {
	commands := &commands{}
	state := &shared.State{
		Flags: commands.getFlags(),
	}

	if len(os.Args) < 2 {
		fmt.Print("Not enough arguments provided.\n")
		os.Exit(1)
	}
	cfg, err := config.Read()
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	state.Cfg = cfg
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
	var loaderProgram *tea.Program

	// Schedule loader start after 500ms
	timer := time.AfterFunc(500*time.Millisecond, func() {
		loaderProgram = tea.NewProgram(output.InitLoader())
		if _, err := loaderProgram.Run(); err != nil {
			log.Printf("Error running loader: %v", err)
		}
	})
	defer timer.Stop()

	result, err := f()

	// Stop the timer if loader hasn't started yet
	if !timer.Stop() && loaderProgram != nil {
		loaderProgram.Send(tea.QuitMsg{})
	}

	return result, err
}

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
		output.ShowResponse(result)
	}
}

func (c *commands) run(s *shared.State, cmd Command) error {
	callFunc := func() (api.Result, error) {
		startTime := time.Now()
		result, err := api.Call(cmd.endpoint, s)
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

func (cmd *commands) getFlags() map[string]any {
	flagsMap := make(map[string]any)

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "--") {
			kv := strings.TrimPrefix(arg, "--")

			parts := strings.SplitN(kv, "=", 2)
			if len(parts) != 2 {
				continue
			}

			val := parts[1]

			// Try to detect type
			if i, err := strconv.Atoi(val); err == nil {
				flagsMap[parts[0]] = i
			} else if b, err := strconv.ParseBool(val); err == nil {
				flagsMap[parts[0]] = b
			} else {
				flagsMap[parts[0]] = val
			}
		}
	}

	return flagsMap
}
