package main

import (
	"os"

	_ "github.com/beshenkaD/maschinenkonzept/admin"
	_ "github.com/beshenkaD/maschinenkonzept/config"
	"github.com/beshenkaD/maschinenkonzept/core"
	_ "github.com/beshenkaD/maschinenkonzept/me"
	_ "github.com/beshenkaD/maschinenkonzept/ping"
)

func main() {
	bot := core.New(os.Getenv("VK_TOKEN"), "/home/beshenka/hueta", true)

	bot.Run()
}
