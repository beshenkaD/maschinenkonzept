package core

import (
	"fmt"
	"strings"
)

const (
	helpCommand = "help"
)

func (b *Bot) help(i *CommandInput) {
	var args string

	for _, arg := range i.Args {
		args += arg + " "
	}

	msg := &Message{
		Text: i.Chat.Prefix + args,
	}

	in := parse(msg, i.Chat, i.User)

	if in == nil {
		b.showAvailableCommands(i.Chat)
		return
	}

	command := commands[in.Command]
	if command == nil {
		b.showAvailableCommands(i.Chat)
		return
	}

	b.showHelp(in, command)
}

func (b *Bot) showHelp(i *CommandInput, help *Command) {
	if help.Description != "" {
		b.sendMessage(i.Chat, fmt.Sprintf(getStrings(i.Chat.Lang).HelpDescription, help.Description))
	}

	var args string
	for _, param := range help.Params {
		args += param.Name + " "
	}

	b.sendMessage(i.Chat, fmt.Sprintf(getStrings(i.Chat.Lang).HelpUsage, i.Chat.Prefix, i.Command, args))
}

func (b *Bot) showAvailableCommands(chat *Chat) {
	var cmds []string

	for k := range commands {
		cmds = append(cmds, k)
	}

	b.sendMessage(chat, fmt.Sprintf(getStrings(chat.Lang).HelpAboutCommand, chat.Prefix))
	b.sendMessage(chat, fmt.Sprintf(getStrings(chat.Lang).HelpAvailableCommands, strings.Join(cmds, ", ")))
}
