package set

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/beshenkaD/maschinenkonzept/core"
)

func set(i *core.CommandInput) (string, error) {
	isAdmin, err := core.IsAdmin(i.Chat, i.User)
	if err != nil {
		return "", err
	}

	if !isAdmin && !i.Chat.IsPrivate() {
		return "", errors.New("вы не являетесь администратором")
	}

	if len(i.Args) < 2 {
		return "", errors.New("недостаточно аргументов")
	}

	switch i.Args[0] {
	case "prefix":
		p := i.Args[1]
		if len(i.Args[1]) > 1 {
			p = p + " "
		}

		i.Chat.Prefix = p
	case "ignore":
		ig, err := strconv.ParseBool(i.Args[1])
		if err != nil {
			return "", fmt.Errorf("неверный аргумент: %s", i.Args[1])
		}

		i.Chat.IgnoreInvalid = ig
	case "command":
		if len(i.Args) < 3 {
			return "", errors.New("вы не указали ни одной команды")
		}

		cmd := func(enable bool) error {
			e := ""
			g := []string{}

			moreThanFive := false
			for j, c := range i.Args[2:] {
				ok := core.IsCommandExist(c)
				if !ok {
					if j < 5 {
						e += fmt.Sprintf("команда `%s` не существует\n", c)
					} else {
						moreThanFive = true
					}
				} else {
					g = append(g, c)
					i.Chat.DisabledCommands[c] = !enable
				}
			}
			if moreThanFive {
				e += fmt.Sprintf("и ещё %d команд", len(i.Args[2:])-5)
			}

			if e != "" {
				if len(g) != 0 {
					var s string
					if enable {
						s = "Включены"
					} else {
						s = "Отключены"
					}
					e += fmt.Sprintf("\n%s команды: %v", s, strings.Join(g, ", "))
				}
				return errors.New(e)
			}

			i.Chat.Save()
			return nil
		}
		switch i.Args[1] {
		case "enable":
			err := cmd(true)
			if err != nil {
				return "", err
			}
		case "disable":
			err := cmd(false)
			if err != nil {
				return "", err
			}
		default:
			return "", fmt.Errorf("неправильный аргумент: %s", i.Args[1])
		}
	default:
		return "", fmt.Errorf("неправильный аргумент: %s. Используйте /help set", i.Args[0])
	}

	i.Chat.Save()

	return "Настройки чата успешно обновлены", nil
}

func get(i *core.CommandInput) (string, error) {
	f := fmt.Sprintf(
		"-- ID чата: %d\n"+
			"-- Префикс: `%s`\n"+
			"-- Игнорировать неправильные команды?: %t\n", i.Chat.ID, i.Chat.Prefix, i.Chat.IgnoreInvalid)

	if len(i.Chat.DisabledCommands) > 0 {
		c := []string{}

		for k, v := range i.Chat.DisabledCommands {
			if v {
				c = append(c, k)
			}
		}

		if len(c) > 0 {
			f += fmt.Sprintf("-- Отключенные команды: %v\n", strings.Join(c, ", "))
		}
	}

	return f, nil
}

func init() {
	core.RegisterCommand(
		"set",
		"Меняет настройки беседы",
		[]core.HelpParam{
			{Name: "prefix", Description: "Префикс для использования команд. Может быть любой длины. Если префикс длинее одного символа то после него необходимо ставить пробел", Optional: true},
			{Name: "ignore", Description: "Игнорировать неправильные команды? Принимает значения `true` или `false`", Optional: true},
			{Name: "command", Description: "Включает или выключает команды для текущего чата. Принимает значения `enable <commands>` или `disable <commands>`", Optional: true},
		},
		set)

	core.RegisterCommand(
		"get",
		"Выводит настройки беседы в удобном виде",
		nil,
		get)
}
