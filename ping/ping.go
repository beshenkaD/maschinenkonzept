package ping

import (
	"github.com/beshenkaD/maschinenkonzept/core"
)

func ping(i *core.CommandInput) (string, error) {
	return "pong", nil
}

func hook1(i *core.HookInput) (string, error) {
	return "a", nil
}

func hook2(i *core.HookInput) (string, error) {
	return "b", nil
}

func init() {
	core.RegisterCommand(
		"ping",
		"ping",
		"alive?",
		[]core.HelpParam{},
		ping)

	core.RegisterHook("a", "hueta", core.ActionPinMessage, hook1)
	core.RegisterHook("b", "hueta", core.ActionPinMessage, hook2)
}
