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
	Session   *api.VK
	lp        *longpoll.LongPoll
	SelfName  string
	SelfID    int
	ChatsLock sync.Mutex
	Chats     map[int]*Chat
	loader    func(*Chat) []Module

	ticker *time.Ticker
	done   chan struct{}

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
		Session:   session,
		lp:        lp,
		SelfName:  group[0].Name,
		SelfID:    group[0].ID,
		Chats:     make(map[int]*Chat),
		loader:    loader,
		Processed: 0,
		StartTime: time.Now(),
	}

	return &b, nil
}

func (b *Bot) AddToConversation(chatID int) {
	conversation := NewConversation(b, chatID)

	b.ChatsLock.Lock()
	b.Chats[chatID] = conversation
	b.ChatsLock.Unlock()

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

func (b *Bot) RunCommand(msg vkMessage, chat *Chat, pm bool) {
	var prefix byte = '/'

	if chat != nil && len(chat.Config.Basic.CommandPrefix) == 1 {
		prefix = chat.Config.Basic.CommandPrefix[0]
	}

	peerID := msg.Message.PeerID
	text := msg.Message.Text

	if len(text) > 1 && text[0] == prefix {
		args := strings.Split(text[1:], " ")
		key := args[0]

		c, ok := chat.commands[key]
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
				go c.Run(msg, args[1:], chat)
			}

		} else if !chat.Config.Basic.IgnoreInvalidCommands {
			vkutil.SendMessage(chat.Bot.Session, "Неправильная команда. Используйте /help", chat.ID, true)
		}
	} else if !pm {
		for _, h := range chat.hooks.OnMessage {
			go h.OnMessage(chat, msg)
		}
	}
}

func (b *Bot) OnMessage(msg vkMessage, chat *Chat) {
	b.Processed++
	pm := chat.ID < 2000000000

	b.RunCommand(msg, chat, pm)
}

func (b *Bot) OnInviteUser(msg vkMessage, chat *Chat) {
	for _, h := range chat.hooks.OnInviteUser {
		if chat.ShouldRunHooks(chat.ID, h) {
			go h.OnInviteUser(chat, msg)
		}
	}
}

func (b *Bot) OnInviteBot(msg vkMessage, chat *Chat) {
	for _, h := range chat.hooks.OnInviteBot {
		if chat.ShouldRunHooks(chat.ID, h) {
			go h.OnInviteBot(chat, msg)
		}
	}
}

func (b *Bot) OnInviteByLink(msg vkMessage, chat *Chat) {
	for _, h := range chat.hooks.OnInviteByLink {
		if chat.ShouldRunHooks(chat.ID, h) {
			go h.OnInviteByLink(chat, msg)
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
		msg := vkMessage(obj)
		peerID := obj.Message.PeerID

		chat, ok := b.Chats[peerID]

		if !ok {
			b.AddToConversation(obj.Message.PeerID)

			chat = b.Chats[peerID]
		}

		switch obj.Message.Action.Type {
		case "chat_photo_update":
		case "chat_photo_remove":
		case "chat_create":
		case "chat_title_update":
		case "chat_invite_user":
			if msg.Message.Action.MemberID == -b.SelfID {
				go b.OnInviteBot(msg, chat)
			} else {
				go b.OnInviteUser(msg, chat)
			}
		case "chat_kick_user":
		case "chat_pin_message":
		case "chat_unpin_message":
		case "chat_invite_user_by_link":
		default:
			go b.OnMessage(msg, chat)
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
