package admin

import (
	"github.com/beshenkaD/maschinenkonzept/core"
)

func init() {
	core.RegisterCommand(
		"kick",
		"Исключает пользователя из беседы. Можно использовать аргументы или ответ на сообщение (только для админов)",
		[]core.HelpParam{
			{Name: "ID(s)", Description: "Один или несколько ID пользователей которых надо исключить", Optional: true},
			{Name: "Упоминания", Description: "Одно или несколько упоминаний пользователей которых надо исключить. Упомянуть пользователя можно с помощью @ или *", Optional: true},
		},
		kick)
}
