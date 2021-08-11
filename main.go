package main

import (
	"log"
	"os"
	"os/signal"

	_ "github.com/beshenkaD/maschinenkonzept/admin"
	"github.com/beshenkaD/maschinenkonzept/core"
	_ "github.com/beshenkaD/maschinenkonzept/me"
	_ "github.com/beshenkaD/maschinenkonzept/ping"
	_ "github.com/beshenkaD/maschinenkonzept/quote"
	_ "github.com/beshenkaD/maschinenkonzept/set"
)

func main() {
	Version := "0.6.3-alpha"
	bot := core.New(os.Getenv("VK_TOKEN"), "/home/beshenka/hueta", Version, true)

	// Handle SIGINT safely
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Println("Safely terminating...")

			if bot != nil {
				bot.Stop()
			} else {
				os.Exit(0)
			}
		}
	}()

	bot.Run()
}
