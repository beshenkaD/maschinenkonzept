package core

import (
	"errors"
	"time"
)

type ResponseHandler func(chatID int, message string)
type ErrorHandler func(chatID int, err error)

type Bot struct {
	ResponseHandler ResponseHandler
	ErrorHandler    ErrorHandler
	Protocol        string  // vk or telegram
	Chats           []*Chat // active chats
	done            chan struct{}
}

func New(h ResponseHandler, e ErrorHandler, protocol string) *Bot {
	b := &Bot{
		ResponseHandler: h,
		ErrorHandler:    e,
		Protocol:        protocol,
		done:            make(chan struct{}),
	}

	go b.startTick()

	return b
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
		for _, chat := range b.Chats {
			for tick := range ticks {
				go b.handleTick(tick, chat)
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func (b *Bot) MessageReceived(chat *Chat, message *Message, sender *User) {
	b.Chats = append(b.Chats, chat)

	input, err := parse(message, chat, sender)

	if err != nil {
		b.SendMessage(chat, err.Error())
		return
	}

	if input == nil {
		go b.handleHook(&HookInput{
			Raw:         message.Text,
			MessageData: message,
			Chat:        chat,
			User:        sender,
		})
		return
	}

	if chat.IsCommandDisabled(input.Command) {
		b.ErrorHandler(chat.ID, errors.New("Command disabled in this chat"))
		return
	}

	switch input.Command {
	case helpCommand:
		go b.help(input)
	default:
		go b.handleCommand(input)
	}
}

func (b *Bot) SendMessage(chat *Chat, message string) {
	if message == "" {
		return
	}

	b.ResponseHandler(chat.ID, message)
}

func (b *Bot) Stop() {
	close(b.done)
}
