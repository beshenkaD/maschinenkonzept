package main

import (
	"github.com/beshenkaD/maschinenkonzept/core"
	"github.com/beshenkaD/maschinenkonzept/core/base"
	"github.com/beshenkaD/maschinenkonzept/core/captcha"
	"log"
	"math/rand"
	"time"
)

const (
	token = "64b80ccddef3594c2b8fb072428241c47e24c78be0f4d07fb818723350ca2b2ee36e0aa55100976affd39"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	modules := make([]core.Module, 0, 2)
	modules = append(modules, base.New())
	modules = append(modules, captcha.New())

	bot, err := core.New(token, '/', modules)

	if err != nil {
		log.Fatal(err)
	}

	bot.Run()
}
