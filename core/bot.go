package core

import (
	"errors"
	"time"
)

type ResponseHandler func(chat int, message string)
type ErrorHandler func(chat int, err error)

type Bot struct {
	ResponseHandler ResponseHandler
	ErrorHandler    ErrorHandler
	Protocol        string // vk or telegram
	Chats           []int  // active chats
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

func (b *Bot) MessageReceived(chat int, message *Message, sender *User) {
	b.Chats = append(b.Chats, chat)

	prefix, ok := prefix[chat]

	if !ok {
		prefix = '/'
	}

	input, err := parse(message, chat, sender, prefix)

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

	if IsCommandDisabled(input.Command, chat) {
		b.ErrorHandler(chat, errors.New("command is disabled"))
		return
	}

	go b.handleCommand(input)
}

func (b *Bot) SendMessage(chat int, message string) {
	if message == "" {
		return
	}

	b.ResponseHandler(chat, message)
}

func (b *Bot) Stop() {
	close(b.done)
}
