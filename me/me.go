package me

import (
	"fmt"
	"strings"

	"github.com/beshenkaD/maschinenkonzept/core"
)

func me(i *core.CommandInput) (string, error) {
	if len(i.Args) == 0 {
		return "", nil
	}

	s := strings.Join(i.Args, " ")
	out := fmt.Sprintf("*%s %s %s", i.User.FirstName, i.User.LastName, s)

	core.DeleteMessages(i.Chat, []int{i.Message.ConversationMessageID})

	return out, nil
}

func init() {
	core.RegisterCommand(
		"me",
		"Позволяет отправить сообщение о себе от третьего лица. Бот также удаляет изначальное сообщение, если может это сделать",
		nil,
		me)
}
