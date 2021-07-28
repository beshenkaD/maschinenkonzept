package core

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Chat struct {
	ID         int
	ConfigLock sync.Mutex
	Config     Config
	hooks      moduleHooks
	Modules    []Module
	commands   map[commandID]Command
	Bot        *Bot
}

// Возвращает объект беседы с дефолтной конфигурацией
func NewChat(bot *Bot, ID int) *Chat {
	return &Chat{
		ID:       ID,
		commands: make(map[commandID]Command),
		Config:   *NewConfig(),
		Bot:      bot,
	}
}

func (ch *Chat) WriteConfig() {
	bs, err := json.Marshal(ch.Config)

	if err != nil {
		fmt.Println(err.Error())
	}

	if err := os.WriteFile(fmt.Sprintf("%d.json", ch.ID), bs, 0664); err != nil {
		fmt.Println(err.Error())
	}
}

func (ch *Chat) LoadConfig() error {
	content, err := os.ReadFile(fmt.Sprintf("./%d.json", ch.ID))

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	var config Config
	err = json.Unmarshal(content, &config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	ch.Config = config

	return nil
}

func (ch *Chat) addCommand(c Command, m Module) {
	name := strings.ToLower(c.Info().Name)
	ch.commands[commandID(name)] = c
}

func (ch *Chat) ShouldRunHooks(m Module) bool {
	if ch.ID < 2000000000 {
		return false
	}

	if _, off := ch.Config.Modules.Disabled[moduleID(strings.ToLower(m.Name()))]; off {
		return false
	}

	return true
}
