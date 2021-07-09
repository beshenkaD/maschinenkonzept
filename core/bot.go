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
	"github.com/beshenkaD/maschinenkonzept/vkutil"
)

type Bot struct {
	Session           *api.VK
	lp                *longpoll.LongPoll
	SelfName          string
	SelfID            int
	ConversationsLock sync.Mutex
	Conversations     map[int]*Conversation
	loader            func(*Conversation) []Module

	ticker *time.Ticker
	done   chan struct{}

	Processed uint
	StartTime time.Time
}

type vkMessage events.MessageNewObject

func New(token string, prefix byte, loader func(*Conversation) []Module) (*Bot, error) {
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
		Session:       session,
		lp:            lp,
		SelfName:      group[0].Name,
		SelfID:        group[0].ID,
		Conversations: make(map[int]*Conversation),
		loader:        loader,
		Processed:     0,
		StartTime:     time.Now(),
	}

	return &b, nil
}

func (b *Bot) AddToConversation(chatID int) {
	conversation := NewConversation(b, chatID)

	b.ConversationsLock.Lock()
	b.Conversations[chatID] = conversation
	b.ConversationsLock.Unlock()

	conversation.Modules = b.loader(conversation)

	for _, v := range conversation.Modules {
		conversation.RegisterModule(v)
		cmds := v.Commands()

		for _, command := range cmds {
			conversation.addCommand(command, v)
		}
	}
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

func (b *Bot) RunCommand(msg events.MessageNewObject, info *Conversation, pm bool) {
	peerID := msg.Message.PeerID
	text := msg.Message.Text

	if len(text) > 1 && text[0] == info.Prefix {
		args := strings.Split(text[1:], " ")
		key := args[0]

		c, ok := info.commands[key]
		if ok {
			if len(args) > 1 {
				if in(args[1], "помощь", "хелп", "help", "usage") {
					vkutil.SendMessage(b.Session, processUsage(c.Usage(), c.Info().Name), peerID, true)
					return
				}
				if in(args[1], "info", "инфо", "информация") {
					vkutil.SendMessage(b.Session, processInfo(c.Info()), peerID, true)
					return
				}
			}

			if pm && !c.Info().ForPm {
				vkutil.SendMessage(b.Session, "Эта команда не работает в лс", peerID, true)
			} else if !pm && !c.Info().ForConf {
				vkutil.SendMessage(b.Session, "Эта команда не работает в беседах", peerID, true)
			} else {
				go c.Run(msg, args[1:], b)
			}

		}
	} else if !pm {
		for _, h := range info.hooks.OnMessage {
			go h.OnMessage(b, msg)
		}
	}
}

func (b *Bot) OnMessage(msg events.MessageNewObject, chat *Conversation) {
	b.Processed++
	pm := chat.ID < 2000000000

	b.RunCommand(msg, chat, pm)
}

func (b *Bot) OnInviteUser(msg events.MessageNewObject, chat *Conversation) {
	for _, h := range chat.hooks.OnInviteUser {
		if chat.ShouldRunHooks(chat.ID, h) {
			go h.OnInviteUser(b, msg)
		}
	}
}

func (b *Bot) OnInviteBot(msg events.MessageNewObject, chat *Conversation) {
	for _, h := range chat.hooks.OnInviteBot {
		if chat.ShouldRunHooks(chat.ID, h) {
			go h.OnInviteBot(b, msg)
		}
	}
}

func (b *Bot) OnInviteByLink(msg events.MessageNewObject, chat *Conversation) {
	for _, h := range chat.hooks.OnInviteByLink {
		if chat.ShouldRunHooks(chat.ID, h) {
			go h.OnInviteByLink(b, msg)
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

func (b *Bot) Run() {
	if b.IsRunning() {
		return
	}
	b.done = make(chan struct{})

	b.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		peerID := obj.Message.PeerID

		chat, ok := b.Conversations[peerID]

		if !ok {
			vkutil.SendMessage(b.Session, "Похоже вы впервые добавили бота в беседу!. Настройте его командой TODO", peerID, true)
			b.AddToConversation(obj.Message.PeerID)

			chat = b.Conversations[peerID]
		}

		switch obj.Message.Action.Type {
		case "chat_photo_update":
		case "chat_photo_remove":
		case "chat_create":
		case "chat_title_update":
		case "chat_invite_user":
			go b.OnInviteUser(obj, chat)
		case "chat_kick_user":
		case "chat_pin_message":
		case "chat_unpin_message":
		case "chat_invite_user_by_link":
		default:
			go b.OnMessage(obj, chat)
		}

		// log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)
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
