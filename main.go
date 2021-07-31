package main

import (
	"os"

	_ "github.com/beshenkaD/maschinenkonzept/ping"
	"github.com/beshenkaD/maschinenkonzept/vk"
)

func main() {
	vk.Run(os.Getenv("VK_TOKEN"), true)
}
