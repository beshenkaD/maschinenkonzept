package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/beshenkaD/maschinenkonzept/core"
)

func config(i *core.CommandInput) (string, error) {
	if len(i.Args) < 2 {
		return "", errors.New("недостаточно аргументов")
	}

	switch strings.ToLower(i.Args[0]) {
	case "prefix", "префикс":
		i.Chat.Prefix = i.Args[1]
		i.Chat.Save()
		return "новый префикс установлен", nil
	case "ignore", "ignoreInvalid", "игнорировать":
		b := true

		switch i.Args[1] {
		case "yes", "true", "да", "y":
			b = true
		case "no", "false", "нет", "n":
			b = false
		default:
			return "", errors.New("неправильные аргументы нах")
		}

		i.Chat.IgnoreInvalid = b
		i.Chat.Save()
		return "параметр установлен", nil
	case "language", "lang", "язык":
		i.Chat.Lang = i.Args[1]
		i.Chat.Save()

		return "параметр установлен", nil
	case "disable", "отключить":
		if len(i.Args) < 3 {
			return "", errors.New("недостаточно аргументов")
		}

		switch i.Args[1] {
		case "command", "команда", "команду":
			for _, arg := range i.Args[2:] {
				if core.IsCommandExist(arg) {
					i.Chat.DisabledCommands[arg] = true
				} else {
					core.SendMessage(i.Chat, fmt.Sprintf("%s command does not exist. Use %shelp to list all commands", arg, i.Chat.Prefix), "", "", nil)
				}
			}
		case "hook", "хук":
			for _, arg := range i.Args[1:] {
				if core.IsHookExist(arg) {
					i.Chat.DisabledHooks[arg] = true
				} else {
					core.SendMessage(i.Chat, fmt.Sprintf("%s hook does not exist. Use %shelp to list all hooks", arg, i.Chat.Prefix), "", "", nil)
				}
			}
		case "tick", "тик":
			for _, arg := range i.Args[1:] {
				if core.IsTickExist(arg) {
					i.Chat.DisabledTicks[arg] = true
				} else {
					core.SendMessage(i.Chat, fmt.Sprintf("%s tick does not exist. Use %shelp to list all ticks", arg, i.Chat.Prefix), "", "", nil)
				}
			}
		}
		i.Chat.Save()
	default:
		return "", errors.New("неправильные аргументы")
	}

	return "", nil
}

func get(i *core.CommandInput) (string, error) {
	f := `
Prefix: %s
Language: %s
Ignore Invalid: %t
Disabled commands: %s
Disabled hooks: %s
Disabled ticks: %s
`
	return fmt.Sprintf(f, i.Chat.Prefix, i.Chat.Lang, i.Chat.IgnoreInvalid, i.Chat.DisabledCommands, i.Chat.DisabledHooks, i.Chat.DisabledTicks), nil
}

func init() {
	core.RegisterCommand("config", "", nil, config)
	core.RegisterCommand("info", "", nil, get)
}
