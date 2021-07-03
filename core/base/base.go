// Базовый модуль бота

package base

import (
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/beshenkaD/maschinenkonzept/apiutil"
	"github.com/beshenkaD/maschinenkonzept/core"
)

type BaseModule struct{}

func New() *BaseModule {
	return &BaseModule{}
}

func (w *BaseModule) Name() string {
	return "Базовый модуль"
}

func (w *BaseModule) OnInviteUser(bot *core.Bot, msg events.MessageNewObject) {
	apiutil.Send(bot.Session, "Привет!", msg.Message.PeerID)
}

func (w *BaseModule) OnKickUser(bot *core.Bot, msg events.MessageNewObject) {
	apiutil.Send(bot.Session, "Пока-пока :(", msg.Message.PeerID)
}

func (w *BaseModule) OnPinMessage(bot *core.Bot, msg events.MessageNewObject) {
	apiutil.Send(bot.Session, "нахуй ты это сделал?", msg.Message.PeerID)
}

func (w *BaseModule) OnInviteBot(bot *core.Bot, msg events.MessageNewObject) {
	apiutil.Send(bot.Session, "Спасибо что добавили меня", msg.Message.PeerID)
}

func (w *BaseModule) OnUnpinMessage(bot *core.Bot, msg events.MessageNewObject) {
	apiutil.Send(bot.Session, "молодец.", msg.Message.PeerID)
}

func (w *BaseModule) Commands() []core.Command {
	return []core.Command{
		&pingCommand{},
	}
}

func (w *BaseModule) Description() string {
	return "Базовый модуль для проверки работоспособности бота"
}

type pingCommand struct{}

func (c *pingCommand) Info() *core.CommandInfo {
	return &core.CommandInfo{
		Name: "Ping",
		Desc: "Проверить работоспособность бота (или поиграть в пинг-понг) :)",
	}
}

func (c *pingCommand) Run(msg events.MessageNewObject, argc int, argv []string, bot *core.Bot) {
	if argc == 0 {
		apiutil.Send(bot.Session, "pong", msg.Message.PeerID)

		return
	}

	if argv[0] == "ru" {
		apiutil.Send(bot.Session, "понг", msg.Message.PeerID)

		return
	}
}

func (c *pingCommand) Usage() *core.CommandUsage {
	return &core.CommandUsage{
		Desc: "Проверяет работспособность бота",
		Params: []core.CommandUsageParam{
			{Name: "ru", Desc: "Бот ответит вам по-русски", Optional: true},
		},
	}
}
