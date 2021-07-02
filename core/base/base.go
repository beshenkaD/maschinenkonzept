// Базовый модуль бота

package base

import (
	"fmt"
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

func (w *BaseModule) OnInviteUser(bot *core.Bot, msg events.MessageNewObject) error {
	apiutil.Send(bot.Session, "Привет!", msg.Message.PeerID)
	return nil
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
		Desc: "Проверить работоспособность бота ( или поиграть в пинг-понг :) )",
	}
}

func (c *pingCommand) Run(msg events.MessageNewObject, args []string, bot *core.Bot) error {
	_, err := apiutil.Send(bot.Session, "pong", msg.Message.PeerID)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (c *pingCommand) Usage() *core.CommandUsage {
	return &core.CommandUsage{}
}
