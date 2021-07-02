package core

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
	"strings"
)

type Bot struct {
	Session  *api.VK
	SelfName string
	SelfID   int
	Prefix   byte
	Modules  []Module
	commands map[string]Command
	hooks    moduleHooks
}

func New(token string, prefix byte, modules []Module) *Bot {
	session := api.NewVK(token)
	group, err := session.GroupsGetByID(nil)

	if err != nil {
		return nil
	}

	b := &Bot{
		Session:  session,
		Prefix:   prefix,
		SelfName: group[0].Name,
		SelfID:   group[0].ID,
		commands: make(map[string]Command),
	}

	for _, m := range modules {
		b.Modules = append(b.Modules, m)
		b.RegisterModule(m)
	}

	return b
}

// TODO: Закончить ветку else, добавить поддержку help'а
func (b *Bot) ProcessCommand(msg events.MessageNewObject) {
	text := msg.Message.Text

	if len(text) > 1 && text[0] == b.Prefix {
		args := strings.Split(text[1:], " ")
		key := args[0]

		c, ok := b.commands[key]
		if ok {
			go c.Run(msg, args[1:], b)
		}
	} else {
		action := msg.Message.Action.Type

		switch action {
		case "chat_invite_user":
			for _, h := range b.hooks.OnInviteUser {
				go h.OnInviteUser(b, msg)
			}
		case "chat_kick_user":
			for _, h := range b.hooks.OnKickUser {
				go h.OnKickUser(b, msg)
			}
		default:
			for _, h := range b.hooks.OnMessage {
				go h.OnMessage(b, msg)
			}
		}
	}
}

func (b *Bot) Run() {
	lp, err := longpoll.NewLongPoll(b.Session, b.SelfID)
	if err != nil {
		log.Fatal(err)
	}

	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)

		b.ProcessCommand(obj)
	})

	log.Println("Start Long Poll")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}
