package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/beshenkaD/maschinenkonzept/core"
	"github.com/beshenkaD/maschinenkonzept/core/base"
	"github.com/beshenkaD/maschinenkonzept/core/captcha"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	modules := make([]core.Module, 0, 2)
	modules = append(modules, base.New())
	modules = append(modules, captcha.New())

	bot, err := core.New(os.Getenv("VK_TOKEN"), '/', modules)

	if err != nil {
		log.Fatal(err)
	}

	bot.Run()
}
