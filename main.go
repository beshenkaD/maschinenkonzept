package main

import (
	"log"
	"os"
	"os/signal"

	_ "github.com/beshenkaD/maschinenkonzept/admin"
	_ "github.com/beshenkaD/maschinenkonzept/config"
	"github.com/beshenkaD/maschinenkonzept/core"
	_ "github.com/beshenkaD/maschinenkonzept/me"
	_ "github.com/beshenkaD/maschinenkonzept/ping"
	_ "github.com/beshenkaD/maschinenkonzept/quote"
)

func main() {
	bot := core.New(os.Getenv("VK_TOKEN"), "/home/beshenka/hueta", "0.5.2-alpha", true)

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
