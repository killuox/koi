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

// Method handlers
func Get(e shared.Endpoint, s *shared.State) (r Result, err error) {
	url := configureUrl(e, s)

	resp, err := http.Get(url)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, err
	}
	return Result{
		Body:   body,
		Url:    url,
		Status: resp.StatusCode,
		Method: e.Method,
	}, nil
}

func Post(e shared.Endpoint, s *shared.State) (r Result, err error) {
	url := configureUrl(e, s)
	payload := map[string]any{}
	for k, p := range e.Parameters {
		val, err := p.GetValue(s, k, e)
		if err != nil && p.Required {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		payload[k] = val
	}

	// Convert map to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Create a POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Body:   body,
		Url:    url,
		Status: resp.StatusCode,
		Method: e.Method,
	}, nil
}

func Update(e shared.Endpoint, s *shared.State) (r Result, err error) {
	url := configureUrl(e, s)
	payload := map[string]any{}
	for k, p := range e.Parameters {
		val, err := p.GetValue(s, k, e)
		if err != nil && p.Required {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		payload[k] = val
	}

	// Convert map to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	req, err := http.NewRequest(e.Method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return Result{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Body:   body,
		Url:    url,
		Status: resp.StatusCode,
		Method: e.Method,
	}, nil
}

func Delete(e shared.Endpoint, s *shared.State) (r Result, err error) {
	url := configureUrl(e, s)

	req, err := http.NewRequest(e.Method, url, nil)
	if err != nil {
		return Result{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, err
	}
	return Result{
		Body:   body,
		Url:    url,
		Status: resp.StatusCode,
		Method: e.Method,
	}, nil
}
