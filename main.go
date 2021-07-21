package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/beshenkaD/maschinenkonzept/core"
)

func loader(c *core.Chat) []core.Module {
	modules := make([]core.Module, 0, 2)
	modules = append(modules, &core.ConfigModule{})
	modules = append(modules, &core.InfoModule{})

	return modules
}

func main() {
	rand.Seed(time.Now().UnixNano())

	bot, err := core.New(os.Getenv("VK_TOKEN"), loader)

	if err != nil {
		log.Fatal(err)
	}

	bot.Run()
}
