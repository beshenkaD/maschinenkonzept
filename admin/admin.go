package admin

import (
	"github.com/beshenkaD/maschinenkonzept/core"
)

func init() {
	core.RegisterCommand(
		"kick",
		"Kick user from conferention",
		[]core.HelpParam{
			{Name: "IDs", Description: "kick by ID(s)", Optional: true},
			{Name: "Mentions", Description: "kick by mention(s)", Optional: true},
			{Name: "Reply", Description: "kick by reply message", Optional: true},
		},
		kick)
}
