package main

import (
	"github.com/killuox/koi/internal/commands"
	"github.com/killuox/koi/internal/env"
)

func main() {
	env.LoadEnv()
	commands.Init()
}
