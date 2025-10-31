package shared

import "github.com/killuox/koi/internal/config"

type State struct {
	Cfg       config.Config
	Flags     map[string]interface{}
	Variables map[string]interface{}
}
