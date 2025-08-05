package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/killuox/koi/internal/config"
)

type Result struct {
	Status   int
	Body     []byte
	Url      string
	Method   string
	Duration time.Duration
}

func Call(e config.Endpoint, cfg config.Config) (r Result, err error) {
	switch e.Method {
	case "GET":
		return Get(e, cfg)
	case "POST":
		return Post(e, cfg)
	// case "PATCH":
	// 	return Patch(e, cfg)
	// case "DELETE":
	// 	return Delete(e, cfg)
	default:
		return Result{}, fmt.Errorf("unsupported method: %s", e.Method)
	}
}

// Utilies
func getUrl(e config.Endpoint, cfg config.Config) string {
	return cfg.API.BaseURL + e.Path
}

// Method handlers
func Get(e config.Endpoint, cfg config.Config) (r Result, err error) {
	//TODO: handle query, path params
	url := getUrl(e, cfg)

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

func Post(e config.Endpoint, cfg config.Config) (r Result, err error) {
	url := getUrl(e, cfg)

	payload := e.Parameters
	// TODO: handle parameters, get defaults value for now
	// payload := map[string]string{
	// 	"name":  "Antonin",
	// 	"email": "antonin@example.com",
	// }

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

// func Patch(e config.Endpoint, cfg config.Config) (r Result, err error)  {}
// func Delete(e config.Endpoint, cfg config.Config) (r Result, err error) {}
