package core

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/beshenkaD/maschinenkonzept/vkutil"
)

type Bot struct {
	Session  *api.VK
	lp       *longpoll.LongPoll
	SelfName string
	SelfID   int
	Prefix   byte
	Modules  []Module
	commands map[string]Command
	hooks    moduleHooks

	ticker *time.Ticker
	done   chan struct{}

	Processed uint
	StartTime time.Time
}

func New(token string, prefix byte, modules []Module) (*Bot, error) {
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
		Prefix:    prefix,
		SelfName:  group[0].Name,
		SelfID:    group[0].ID,
		commands:  make(map[string]Command),
		Processed: 0,
		StartTime: time.Now(),
		Modules:   modules,
	}

	for _, m := range b.Modules {
		b.RegisterModule(m)
	}

	return &b, nil
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
	s := "⚙ %s -- %s"

	return fmt.Sprintf(s, info.Name, info.Desc)
}

func (b *Bot) ProcessMessage(msg events.MessageNewObject) {
	b.Processed++

	peerID := msg.Message.PeerID
	text := msg.Message.Text

	if len(text) > 1 && text[0] == b.Prefix {
		args := strings.Split(text[1:], " ")
		key := args[0]

		c, ok := b.commands[key]
		if ok {
			if len(args) > 1 {
				if in(args[1], "помощь", "хелп", "help", "usage") {
					_, err := vkutil.SendMessage(b.Session, processUsage(c.Usage(), c.Info().Name), peerID, true)
					if err != nil {
						log.Println(err.Error(), "peer_id: ", peerID)
					}
					return
				}
				if in(args[1], "info", "инфо", "информация") {
					_, err := vkutil.SendMessage(b.Session, processInfo(c.Info()), peerID, true)
					if err != nil {
						log.Println(err.Error(), "peer_id: ", peerID)
					}
					return
				}
			}

			go c.Run(msg, args[1:], b)

			for _, h := range b.hooks.OnCommand {
				go h.OnCommand(b, msg)
			}
		}
	} else {
		switch msg.Message.Action.Type {
		case "chat_invite_user":
			if msg.Message.Action.MemberID == -b.SelfID {
				for _, h := range b.hooks.OnInviteBot {
					go h.OnInviteBot(b, msg)
				}
			} else {
				for _, h := range b.hooks.OnInviteUser {
					go h.OnInviteUser(b, msg)
				}
			}
		case "chat_kick_user":
			for _, h := range b.hooks.OnKickUser {
				go h.OnKickUser(b, msg)
			}
		case "chat_pin_message":
			for _, h := range b.hooks.OnPinMessage {
				go h.OnPinMessage(b, msg)
			}
		case "chat_unpin_message":
			for _, h := range b.hooks.OnUnpinMessage {
				go h.OnUnpinMessage(b, msg)
			}
		case "chat_invite_user_by_link":
			for _, h := range b.hooks.OnInviteByLink {
				go h.OnInviteByLink(b, msg)
			}
		case "chat_create":
			for _, h := range b.hooks.OnChatCreate {
				go h.OnChatCreate(b, msg)
			}
		default:
			for _, h := range b.hooks.OnMessage {
				go h.OnMessage(b, msg)
			}
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

	go func() {
		b.ticker = time.NewTicker(time.Second)
		for {
			select {
			case <-b.done:
				b.ticker.Stop()
				return
			case <-b.ticker.C:
			}
			for _, h := range b.hooks.OnTick {
				go h.OnTick(b)
			}
		}
	}()

	b.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		// log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)

		b.ProcessMessage(obj)
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
