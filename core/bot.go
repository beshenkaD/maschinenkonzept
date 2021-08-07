package core

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

var Vk *api.VK

type Bot struct {
	chats map[int]*Chat // active chats
	done  chan struct{}
}

// TODO: translation support
func New(token string, debug bool) *Bot {
	Vk = api.NewVK(token)

	return &Bot{
		chats: make(map[int]*Chat),
		done:  make(chan struct{}),
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
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

func (b *Bot) startTick() {
	for b.IsRunning() {
		for _, chat := range b.chats {
			for tick := range ticks {
				go b.handleTick(tick, chat)
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func (b *Bot) messageReceived(chat *Chat, message *Message, sender *User) {
	b.chats[chat.ID] = chat

	input, err := parse(message, chat, sender)

	if err != nil {
		b.sendMessage(chat, err.Error())
		return
	}

	if input == nil {
		go b.handleHook(&HookInput{
			Message: message,
			Chat:    chat,
			User:    sender,
		})
		return
	}

	if chat.IsCommandDisabled(input.Command) {
		b.sendError(chat, errors.New("Command disabled in this chat"))
		return
	}

	switch input.Command {
	case helpCommand:
		go b.help(input)
	default:
		go b.handleCommand(input)
	}
}

func (b *Bot) sendMessage(chat *Chat, message string) {
	if message == "" {
		return
	}

	bu := params.NewMessagesSendBuilder()
	bu.PeerID(chat.ID)
	bu.RandomID(0)
	bu.Message(message)

	_, err := Vk.MessagesSend(bu.Params)

	if err != nil {
		log.Println(err.Error())
	}
}

func (b *Bot) sendError(chat *Chat, err error) {
	if err == nil {
		return
	}

	bu := params.NewMessagesSendBuilder()
	bu.PeerID(chat.ID)
	bu.RandomID(0)
	bu.Message("Ошибка: " + err.Error())

	_, e := Vk.MessagesSend(bu.Params)

	if e != nil {
		log.Println(err.Error())
	}
}

func (b *Bot) getUser(ID int) *User {
	if ID < 0 {
		bu := params.NewGroupsGetByIDBuilder()
		bu.GroupID(strconv.Itoa(ID))
		bu.Lang(0)

		bot, err := Vk.GroupsGetByID(bu.Params)

		if err != nil {
			return &User{
				ID:        ID,
				FirstName: "",
				LastName:  "",
				IsBot:     true,
			}
		}

		return &User{
			ID:        ID,
			FirstName: bot[0].Name,
			LastName:  "",
			IsBot:     true,
		}
	}

	bu := params.NewUsersGetBuilder()
	bu.Lang(0)
	bu.UserIDs([]string{strconv.Itoa(ID)})

	users, err := Vk.UsersGet(bu.Params)

	if err != nil {
		return &User{
			ID:        ID,
			FirstName: "",
			LastName:  "",
			IsBot:     false,
		}
	}

	return &User{
		ID:        ID,
		FirstName: users[0].FirstName,
		LastName:  users[0].LastName,
		IsBot:     false,
	}
}

func (b *Bot) Run() {
	group, err := Vk.GroupsGetByID(nil)

	if err != nil {
		log.Fatal(err)
	}

	lp, err := longpoll.NewLongPoll(Vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	go b.startTick()

	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		chat, ok := b.chats[obj.Message.PeerID]

		if !ok {
			chat = newChat(obj.Message.PeerID)
		}

		message := Message(obj.Message)
		sender := b.getUser(obj.Message.FromID)

		b.messageReceived(chat, &message, sender)
	})

	log.Println("Start Long Poll (VK)")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}

func (b *Bot) Stop() {
	close(b.done)
}
