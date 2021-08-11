package admin

import (
	"errors"
	"strconv"
	"strings"

	"github.com/beshenkaD/maschinenkonzept/core"
)

func parseMention(s string) int {
	var ID = s
	var bot bool

	if strings.HasPrefix(s, "[id") {
		ID = strings.TrimPrefix(s, "[id")
		bot = false
	} else {
		ID = strings.TrimPrefix(ID, "[club")
		bot = true
	}

	ID = strings.TrimSuffix(ID, "]")
	ID = strings.Split(ID, "|")[0]

	sp, _ := strconv.Atoi(ID)

	if bot {
		sp = -sp
	}

	return sp
}

func kick(i *core.CommandInput) (string, error) {
	members, err := core.GetConversationMembers(i.Chat)

	if err != nil {
		return "", err
	}

	found := false
	for _, m := range members {
		if m.IsAdmin && m.MemberID == i.User.ID {
			found = true
			break
		}
	}

	if !found {
		return "", errors.New("вы не являетесь администратором")
	}

	isMention := func(s string) bool {
		return (strings.HasPrefix(s, "[id") || strings.HasPrefix(s, "[club")) && strings.HasSuffix(s, "]")
	}

	if len(i.Args) > 0 {
		for _, arg := range i.Args {
			if isMention(arg) {
				if err := core.RemoveUser(i.Chat, parseMention(arg)); err != nil {
					return "", err
				}
			} else if ID, err := strconv.Atoi(arg); err == nil {
				if err := core.RemoveUser(i.Chat, ID); err != nil {
					return "", err
				}
			} else {
				return "", errors.New("Неправильный аргумент: " + `"` + arg + `"`)
			}
		}
		return "", nil
	}

	if i.Message.ReplyMessage != nil {
		ID := i.Message.ReplyMessage.FromID
		if err := core.RemoveUser(i.Chat, ID); err != nil {
			return "", err
		}
	}

	return "", errors.New("вы не передали никаких аргументов и не ответили ни на какое сообщение")
}
