package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/killuox/koi/internal/config"
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

type UrlConfig struct {
	Url string
}

var validMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

func Call(e config.Endpoint, s *shared.State) (r Result, err error) {
	if !slices.Contains(validMethods, e.Method) {
		return Result{}, fmt.Errorf("invalid method: %s", e.Method)
	} else {
		url := configureUrl(e, s)
		return doRequest(url, e, s)
	}
}

func doRequest(url string, e config.Endpoint, s *shared.State) (Result, error) {
	var body io.Reader

	// Get parameters values for payload
	if e.Method == http.MethodPost || e.Method == http.MethodPut || e.Method == http.MethodPatch {
		payload := map[string]any{}
		setPayload(e, payload, s)
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

func setPayload(e config.Endpoint, p map[string]any, s *shared.State) {
	for k, param := range e.Parameters {
		val, err := param.GetValue(s.Flags, k, e)
		if err != nil && param.Required {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		p[k] = val
	}
}

func configureUrl(e config.Endpoint, s *shared.State) string {
	path := e.Path
	var queryCount int
	var re = regexp.MustCompile(`\{[^}]+\}`) // Check if {anything}
	if re.MatchString(path) {
		// Check if we need to replace a dynamic path parameters
		for k, p := range e.Parameters {
			val, err := p.GetValue(s.Flags, k, e)
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
