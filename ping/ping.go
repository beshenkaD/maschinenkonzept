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
		"ping",
		"Check bot life",
		[]core.HelpParam{},
		ping)
}
