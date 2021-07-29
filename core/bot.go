package core

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

type Bot struct {
	Session     *api.VK
	lp          *longpoll.LongPoll
	SelfName    string
	SelfID      int
	ChatsLock   sync.RWMutex
	Chats       map[int]*Chat
	ConfigsPath string
	loader      func(*Chat) []Module

	done chan struct{}

	Processed uint
	StartTime time.Time
}

func New(token string, loader func(*Chat) []Module) (*Bot, error) {
	session := api.NewVK(token)
	group, err := session.GroupsGetByID(nil)

	if err != nil {
		return nil, err
	}

	lp, err := longpoll.NewLongPoll(session, group[0].ID)

	if err != nil {
		return nil, err
	}

	b := Bot{
		Session:     session,
		lp:          lp,
		SelfName:    group[0].Name,
		SelfID:      group[0].ID,
		Chats:       make(map[int]*Chat),
		ConfigsPath: "/home/beshenka/go/src/github.com/beshenkaD/maschinenkonzept/res",
		loader:      loader,
		Processed:   0,
		StartTime:   time.Now(),
	}

	return &b, nil
}

func (b *Bot) AddChat(chatID int) {
	chat := NewChat(b, chatID)

	b.ChatsLock.Lock()
	b.Chats[chatID] = chat
	b.ChatsLock.Unlock()

	chat.Modules = b.loader(chat)

	for _, v := range chat.Modules {
		chat.RegisterModule(v)
		cmds := v.Commands()

		for _, command := range cmds {
			chat.addCommand(command, v)
		}
	}

	chat.LoadConfig()
}

func processUsage(usage *CommandUsage, name string) string {
	s := "📝%s -- %s\n\n"

	s = fmt.Sprintf(s, strings.ToLower(name), usage.Desc)

	if len(usage.Params) > 0 {
		i := 1
		s += "⚙ Параметры:\n"

		for _, p := range usage.Params {
			var opt string

			if p.Optional {
				opt = ""
			} else {
				opt = "(обязательный)"
			}

			s += fmt.Sprintf("%d. %s -- %s %s\n", i, p.Name, p.Desc, opt)

			i++
		}
	}

	return s
}

func processInfo(info *CommandInfo) string {
	s := "⚙ %s -- %s\n"

	if info.ForConf && !info.ForPm {
		s += "Только для бесед 🍹"
	}

	if info.ForPm && !info.ForConf {
		s += "Только для личных сообщений бота 🔖"
	}

	return fmt.Sprintf(s, info.Name, info.Desc)
}

func runCommand(msg vkMessage, chat *Chat, pm bool) {
	var prefix byte = '/'

	if chat != nil && len(chat.Config.Basic.CommandPrefix) == 1 {
		prefix = chat.Config.Basic.CommandPrefix[0]
	}

	text := msg.Message.Text

	if len(text) > 1 && text[0] == prefix {
		args := strings.Split(text[1:], " ")
		key := commandID(args[0])

		c, ok := chat.commands[key]
		if !ok {
			if alias, aliasok := chat.Config.Basic.Aliases[string(key)]; aliasok {
				var m string

				if len(args) > 1 {
					m = chat.Config.Basic.CommandPrefix + alias + " "

					for _, a := range args[1:] {
						m += a
					}
				} else {
					m = chat.Config.Basic.CommandPrefix + alias
				}

				args = strings.Split(m[1:], " ")

				if len(args) < 1 {
					chat.SendMessage("Ты ахуел?")
					return
				}

				arg := commandID(strings.ToLower(args[0]))
				c, ok = chat.commands[arg]
			}
		}
		if ok {
			if len(args) > 1 {
				if in(args[1], "помощь", "хелп", "help", "usage") {
					chat.SendMessage(processUsage(c.Usage(), c.Info().Name))
					return
				}
				if in(args[1], "info", "инфо", "информация") {
					chat.SendMessage(processInfo(c.Info()))
					return
				}
			}

			if disabled := chat.Config.Modules.DisabledCommands[key]; disabled {
				chat.SendMessage("Эта команда отключена в данной беседе")
				return
			}

			if pm && !c.Info().ForPm {
				chat.SendMessage("Эта команда не работает в лс")
			} else if !pm && !c.Info().ForConf {
				chat.SendMessage("Эта команда не работает в беседах")
			} else {
				out := c.Run(msg, args[1:], chat)
				chat.SendMessage(out)
			}

		} else if !chat.Config.Basic.IgnoreInvalidCommands {
			chat.SendMessage("Неправильная команда. Используйте help")
		}
	} else if !pm {
		for _, h := range chat.hooks.OnMessage {
			go h.OnMessage(chat, msg)
		}
	}
}

func (b *Bot) IsRunning() bool {
	if b.done != nil {
		select {
		case <-b.done:
		default:
			return true
		}
	}
	return false
}

func (b *Bot) OnTick() {
	time.Sleep(2 * time.Second)

	for b.IsRunning() {
		b.ChatsLock.RLock()

		chats := make([]*Chat, 0, len(b.Chats))
		for _, v := range b.Chats {
			chats = append(chats, v)
		}

		b.ChatsLock.RUnlock()

		for _, chat := range chats {
			for _, h := range chat.hooks.OnTick {
				if chat.ShouldRunHooks(h) {
					h.OnTick(chat)
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (b *Bot) Run() {
	if b.IsRunning() {
		return
	}
	b.done = make(chan struct{})

	go b.OnTick()

	b.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		msg := vkMessage(obj)
		peerID := obj.Message.PeerID

		chat, ok := b.Chats[peerID]

		if !ok {
			b.AddChat(obj.Message.PeerID)

			chat = b.Chats[peerID]
		}

		switch obj.Message.Action.Type {
		case "chat_photo_update":
			go b.OnPhotoRemove(chat, msg)
		case "chat_photo_remove":
			go b.OnPhotoRemove(chat, msg)
		case "chat_create":
			go b.OnChatCreate(chat, msg)
		case "chat_title_update":
			go b.OnTitleUpdate(chat, msg)
		case "chat_invite_user":
			go b.OnInviteUser(chat, msg)
		case "chat_kick_user":
			go b.OnKickUser(chat, msg)
		case "chat_pin_message":
			go b.OnPinMessage(chat, msg)
		case "chat_unpin_message":
			go b.OnUnpinMessage(chat, msg)
		case "chat_invite_user_by_link":
			go b.OnInviteByLink(chat, msg)
		default:
			go b.OnMessage(chat, msg)
		}
	})

	log.Println("Start Long Poll")

	if err := b.lp.Run(); err != nil {
		log.Println(err.Error())
	}
}

func (b *Bot) Stop() {
	if !b.IsRunning() {
		return
	}
	b.lp.Shutdown()
	close(b.done)
}
