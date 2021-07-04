package core

import (
	"context"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/beshenkaD/maschinenkonzept/vkutil"
	"log"
	"strings"
	"time"
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
	working  bool

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

	b := &Bot{
		Session:   session,
		lp:        lp,
		Prefix:    prefix,
		SelfName:  group[0].Name,
		SelfID:    group[0].ID,
		commands:  make(map[string]Command),
		Processed: 0,
		StartTime: time.Now(),
		working:   true,
	}

	for _, m := range modules {
		b.Modules = append(b.Modules, m)
		b.RegisterModule(m)
	}

	return b, nil
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

func in(s string, a []string) bool {
	for _, e := range a {
		if e == s {
			return true
		}
	}

	return false
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
				if in(args[1], []string{"помощь", "хелп", "help", "usage"}) {
					_, err := vkutil.SendMessage(b.Session, processUsage(c.Usage(), c.Info().Name), peerID, true)
					if err != nil {
						log.Println(err.Error(), "peer_id: ", peerID)
					}
					return
				}
				if in(args[1], []string{"info", "инфо", "информация"}) {
					_, err := vkutil.SendMessage(b.Session, processInfo(c.Info()), peerID, true)
					if err != nil {
						log.Println(err.Error(), "peer_id: ", peerID)
					}
					return
				}
			}

			go c.Run(msg, len(args[1:]), args[1:], b)

			for _, h := range b.hooks.OnCommand {
				go h.OnCommand(b, msg)
			}
		}
	} else {
		action := msg.Message.Action.Type

		switch action {
		case "chat_invite_user":
			if msg.Message.Action.MemberID == (b.SelfID * -1) {
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

// Запускает бота. Если хочешь иметь возможность его отключить потом то запускай в горутине
func (b *Bot) Run() {
	go func() {
		if !b.working {
			return
		}

		for range time.Tick(time.Second) {
			for _, h := range b.hooks.OnTick {
				go h.OnTick(b)
			}
		}
	}()

	b.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)

		b.ProcessMessage(obj)
	})

	log.Println("Start Long Poll")

	// Это довольно хуевый расклад, надо бы переделать на каналы/ещё как-нибудь
	// Чтоб можно было ловить ошибки отсюда запуская в горутине. Иначе метод Stop() просто не будет ничего делать
	if err := b.lp.Run(); err != nil {
		log.Println(err.Error())
	}
}

func (b *Bot) Stop() {
	b.lp.Shutdown()
	b.working = false
}
