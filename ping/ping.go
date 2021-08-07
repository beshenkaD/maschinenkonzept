package ping

import (
	"github.com/beshenkaD/maschinenkonzept/core"
)

func ping(i *core.CommandInput) (string, error) {
	return "pong", nil
}

func init() {
	core.RegisterCommand(
		"ping",
		"alive?",
		nil,
		ping)
}
