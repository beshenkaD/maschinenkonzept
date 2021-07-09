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
func NewConversation(bot *Bot, ID int) *Chat {
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

func (ch *Chat) ShouldRunHooks(chatID int, m Module) bool {
	if chatID < 2000000000 {
		return false
	}

	if _, off := ch.Config.Modules.Disabled[strings.ToLower(m.Name())]; off {
		return false
	}

	return true
}
