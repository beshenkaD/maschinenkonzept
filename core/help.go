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

	s := fmt.Sprintf("%s version %s\n\nВведите %shelp <command> чтобы получить детальное описание команды.\n\n", b.Name, b.Version, i.Chat.Prefix)
	s += getAvailableCommands(i.Chat)

	if in == nil {
		b.sendMessage(i.Chat, s)
		return
	}

	command := commands[in.Command]
	if command == nil {
		b.sendMessage(i.Chat, s)
		return
	}

	b.sendMessage(i.Chat, getHelp(in, command))
}

func getHelp(i *CommandInput, help *Command) string {
	s := fmt.Sprintf("%s: %s\n", help.Trigger, help.Description)

	if len(help.Params) != 0 {
		s += "Аргументы:\n"
	}

	for _, param := range help.Params {
		optional := "[Обязательный]"

		if param.Optional {
			optional = ""
		}

		s += fmt.Sprintf("-- %s: %s %s\n", param.Name, param.Description, optional)
	}

	return s
}

func getAvailableCommands(chat *Chat) string {
	var cmds []string

	for k := range commands {
		if chat.IsCommandDisabled(k) {
			k += " [Отключена]"
		}
		cmds = append(cmds, k)
	}

	return fmt.Sprintf("Доступные команды: %v", strings.Join(cmds, ", "))
}
