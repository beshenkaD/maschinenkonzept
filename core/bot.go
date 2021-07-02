package core

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
)

// TODO: добавить помимо command ещё и listener поверх которого можно будет
// Реализовать антиспам/сообщения по таймеру/реакцию на вход/выход из беседы и баны

// listener - команда которая реагирует на все сообщения
type command func(session *api.VK, message events.MessageNewObject)

type Bot struct {
	Session   *api.VK
	Prefix    string
	Commands  map[string]command
	Listeners []command
}

func (b *Bot) RunListeners(message events.MessageNewObject) {
	for _, l := range b.Listeners {
		go l(b.Session, message)
	}
}

func (b *Bot) RegisterCommand(name string, proc command) {
	b.Commands[b.Prefix+name] = proc
}

func (b *Bot) UnregisterCommand(name string) {
	delete(b.Commands, b.Prefix+name)
}

func NewBot(token, prefix string) *Bot {
	commands := make(map[string]command)
	commands[prefix+"ping"] = ping
	commands[prefix+"stat"] = stat

	listeners := make([]command, 1)
	listeners[0] = hello

	return &Bot{
		Session:   api.NewVK(token),
		Prefix:    prefix,
		Commands:  commands,
		Listeners: listeners,
	}
}
