package ping

import (
	"fmt"

	"github.com/beshenkaD/maschinenkonzept/core"
)

func info(i *core.CommandInput) (string, error) {
	out := fmt.Sprintf("ChatID: %d\nUser: [id%d|%s %s]\nSettings:\nIgnoreInvalid: %t\nPreifx: %s\nDisabled Commands: %v\nHooks: %v\nTicks: %v",
		i.Chat.ID,
		i.User.ID,
		i.User.FirstName,
		i.User.LastName,
		i.Chat.IgnoreInvalid,
		i.Chat.Prefix,
		i.Chat.DisabledCommands,
		i.Chat.DisabledHooks,
		i.Chat.DisabledTicks)
	return out, nil
}

func init() {
	core.RegisterCommand(
		"info",
		"info",
		"Information about current chat",
		[]core.HelpParam{},
		info)
}
