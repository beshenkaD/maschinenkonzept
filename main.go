package main

import (
	"github.com/beshenkaD/maschinenkonzept/core"
	"github.com/beshenkaD/maschinenkonzept/plugins/admin"
)

func main() {
	bot := core.NewBot("64b80ccddef3594c2b8fb072428241c47e24c78be0f4d07fb818723350ca2b2ee36e0aa55100976affd39", "/")
	bot.RegisterCommand("kick", admin.Kick)
	bot.Run()
}
