package core

import (
	"strings"
	"sync"
)

type Chat struct {
	ID         int
	ConfigLock sync.Mutex
	Config     Config
	hooks      moduleHooks
	Modules    []Module
	commands   map[string]Command
	Bot        *Bot
}

// Возвращает объект беседы с дефолтной конфигурацией
func NewChat(bot *Bot, ID int) *Chat {
	return &Chat{
		ID:       ID,
		commands: make(map[string]Command),
		Config:   *DefaultConfig(),
		Bot:      bot,
	}
}

func (ch *Chat) addCommand(c Command, m Module) {
	name := strings.ToLower(c.Info().Name)
	ch.commands[name] = c
}

func (ch *Chat) ShouldRunHooks(m Module) bool {
	if ch.ID < 2000000000 {
		return false
	}

	if _, off := ch.Config.Modules.Disabled[strings.ToLower(m.Name())]; off {
		return false
	}

	return true
}
