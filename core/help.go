package core

import (
	"fmt"
	"strings"
)

const (
	helpDescripton    = "Description: %s"
	helpUsage         = "Usage: %s%s %s"
	availableCommands = "Available commands: %v"
	helpAboutCommand  = "Type: '%shelp <command>' to see details about a specific command."
	helpCommand       = "help"
)

func (b *Bot) help(i *CommandInput) {
	var args string

	for _, arg := range i.Args {
		args += arg + " "
	}

	msg := &Message{
		Text: string(getPrefix(i.Chat)) + args,
	}

	in, _ := parse(msg, i.Chat, i.User, getPrefix(i.Chat))

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
		b.SendMessage(i.Chat, fmt.Sprintf(helpDescripton, help.Description))
	}

	var args string
	for _, param := range help.Params {
		args += param.Name + " "
	}

	b.SendMessage(i.Chat, fmt.Sprintf(helpUsage, getPrefix(i.Chat), i.Command, args))
}

func (b *Bot) showAvailableCommands(chat int) {
	var cmds []string

	for k := range commands {
		cmds = append(cmds, k)
	}

	b.SendMessage(chat, fmt.Sprintf(helpAboutCommand, getPrefix(chat)))
	b.SendMessage(chat, fmt.Sprintf(availableCommands, strings.Join(cmds, ", ")))
}
