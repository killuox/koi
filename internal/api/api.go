package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/killuox/koi/internal/shared"
	"github.com/killuox/koi/internal/utils"
	"github.com/killuox/koi/internal/variables"
)

type Result struct {
	Status   int
	Body     []byte
	Url      string
	Method   string
	Duration time.Duration
}

func Call(e shared.Endpoint, s *shared.State) (r Result, err error) {
	switch e.Method {
	case "GET":
		return Get(e, s)
	case "POST":
		return Post(e, s)
	case "PUT", "PATCH":
		return Update(e, s)
	case "DELETE":
		return Delete(e, s)
	default:
		return Result{}, fmt.Errorf("unsupported method: %s", e.Method)
	}
}

// Utilities
func configureUrl(e shared.Endpoint, s *shared.State) string {
	path := e.Path
	var queryCount int
	var re = regexp.MustCompile(`\{[^}]+\}`) // Check if {anything}
	if re.MatchString(path) {
		// Check if we need to replace a dynamic path parameters
		for k, p := range e.Parameters {
			val, err := p.GetValue(s, k, e)
			if err != nil && p.Required {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}
			switch p.In {
			case "path":
				path = strings.ReplaceAll(path, fmt.Sprintf("{%s}", k), fmt.Sprintf("%v", val))
			case "query":
				sep := "?"
				if queryCount > 0 {
					sep = "&"
				}
				path += fmt.Sprintf("%s%s=%s", sep, k, val)
			}
		}
	}
	return s.Cfg.API.BaseURL + path
}

func Get(e shared.Endpoint, s *shared.State) (Result, error)    { return doRequest(e, s) }
func Post(e shared.Endpoint, s *shared.State) (Result, error)   { return doRequest(e, s) }
func Update(e shared.Endpoint, s *shared.State) (Result, error) { return doRequest(e, s) }
func Delete(e shared.Endpoint, s *shared.State) (Result, error) { return doRequest(e, s) }

func doRequest(e shared.Endpoint, s *shared.State) (Result, error) {
	url := configureUrl(e, s)

	// Build payload only for methods that can have a body
	var body io.Reader
	if e.Method == http.MethodPost || e.Method == http.MethodPut || e.Method == http.MethodPatch {
		payload := map[string]any{}
		for k, p := range e.Parameters {
			val, err := p.GetValue(s, k, e)
			if err != nil && p.Required {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}
			payload[k] = val
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return Result{}, fmt.Errorf("error encoding JSON: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	// Build request
	req, err := http.NewRequest(e.Method, url, body)
	if err != nil {
		return Result{}, err
	}

	// Set headers
	for key, val := range s.Cfg.API.Headers {
		req.Header.Set(key, val)
	}

	// Send
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, err
	}

	if e.SetVariables.Body != nil {
		// Unmarshal response body into a generic map
		var respMap map[string]any
		if err := json.Unmarshal(respBody, &respMap); err != nil {
			return Result{}, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		for varName, sourceName := range e.SetVariables.Body {
			source, ok := sourceName.(string)
			if !ok {
				continue // skip invalid config
			}

			// Navigate response JSON using dot notation path
			if val, found := utils.DeepGet(respMap, source); found {
				variables.SetUserVariable(varName, val)
			}
		}
	}

	return Result{
		Body:   respBody,
		Url:    url,
		Status: resp.StatusCode,
		Method: e.Method,
	}, nil
}
