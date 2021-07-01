package core

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
)

// TODO: добавить помимо command ещё и listener поверх которого можно будет
// Реализовать антиспам/сообщения по таймеру/реакцию на вход/выход из беседы и баны

type command func(session *api.VK, message events.MessageNewObject)

type Bot struct {
	Token    string
	Prefix   string
	Commands map[string]command
}

func (b *Bot) Register(name string, proc command) {
    b.Commands[name] = proc
}

func (b *Bot) Unregister(name string) {
    delete(b.Commands, name)
}

func NewBot(token, prefix string) Bot {
    commands := make(map[string]command)
	return Bot{
		Token:  token,
		Prefix: prefix,
        Commands: commands,
	}
}
