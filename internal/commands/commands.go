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
	"github.com/go-playground/validator/v10"
	"github.com/killuox/koi/internal/api"
	"github.com/killuox/koi/internal/config"
	"github.com/killuox/koi/internal/output"
	"github.com/killuox/koi/internal/shared"
	"github.com/killuox/koi/internal/variables"
)

type Command struct {
	name      string
	args      []string
	endpoint  config.Endpoint
	variables map[string]interface{}
}

type Cli struct{}

func Init() {
	cli := &Cli{}

	vars, err := variables.GetUserVariables()
	if err != nil {
		fmt.Print("Error while getting user variables")
		os.Exit(1)
	}

	cfg := config.Config{}

	err = cfg.Init(vars)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if err := cfg.Validate(cfg); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			fmt.Println("❌ Invalid koi.config.yaml:")
			for _, e := range ve {
				fmt.Printf("  - %s: %s\n", e.Namespace(), cfg.CreateValidatorMessage(e))
			}
		} else {
			fmt.Printf("❌ Config error: %s\n", err)
		}
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		cli.printHelp(cfg)
		return
	}

	state := &shared.State{
		Flags:     cli.getFlags(),
		Cfg:       cfg,
		Variables: vars,
	}

	cName := os.Args[1]
	args := os.Args[2:]

	ep, ok := cfg.Endpoints[cName]
	if !ok {
		fmt.Printf("no endpoints found for %s\n", cName)
		os.Exit(1)
	}

	cmd := Command{
		name:      cName,
		args:      args,
		endpoint:  ep,
		variables: vars,
	}

	err = cli.run(state, cmd)
	if err != nil {
		fmt.Printf("Error while running the command: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func (c *Cli) runWithLoader(
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

	result, err := f()

	// Stop the timer if loader hasn't started yet
	if !timer.Stop() && loaderProgram != nil {
		loaderProgram.Send(tea.QuitMsg{})
	}

	return result, err
}

func (c *Cli) processAPIResult(
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

func (c *Cli) run(s *shared.State, cmd Command) error {
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

func (cmd *Cli) getFlags() map[string]any {
	flagsMap := make(map[string]any)

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Long flags: --key=value or --key value
		if strings.HasPrefix(arg, "--") {
			kv := strings.TrimPrefix(arg, "--")

			if strings.Contains(kv, "=") {
				parts := strings.SplitN(kv, "=", 2)
				flagsMap[parts[0]] = parseValue(parts[1])
			} else {
				// If next arg exists and isn't a flag, use it as value
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					flagsMap[kv] = parseValue(args[i+1])
					i++
				} else {
					flagsMap[kv] = true
				}
			}

			// Short flags: -k value or -k=value
		} else if strings.HasPrefix(arg, "-") {
			kv := strings.TrimPrefix(arg, "-")

			if strings.Contains(kv, "=") {
				parts := strings.SplitN(kv, "=", 2)
				flagsMap[parts[0]] = parseValue(parts[1])
			} else {
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					flagsMap[kv] = parseValue(args[i+1])
					i++
				} else {
					flagsMap[kv] = true
				}
			}
		}
	}

	return flagsMap
}

func (cmd *Cli) printHelp(cfg config.Config) {
	fmt.Println("koi - API Testing CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  koi <endpoint> [options]")
	fmt.Println()
	fmt.Println("Available Endpoints:")

	for name, ep := range cfg.Endpoints {
		fmt.Printf("  %-12s %s %s\n", name, ep.Method, ep.Path)
	}

	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  koi login --email=user@example.com --password=secret")
	fmt.Println("  koi health")
	fmt.Println()
	fmt.Println("Use \"koi help <endpoint>\" for more information about an endpoint.")
}

func parseValue(val string) any {
	// Try int
	if i, err := strconv.Atoi(val); err == nil {
		return i
	}
	// Try float
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	}
	// Try bool
	if b, err := strconv.ParseBool(val); err == nil {
		return b
	}
	// Default: string
	return val
}
