package quote

import (
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/beshenkaD/maschinenkonzept/core"
)

type QuoteModule struct{}

func New() *QuoteModule {
	return &QuoteModule{}
}

func (w *QuoteModule) Name() string {
	return "Генератор цитат"
}

func (w *QuoteModule) Description() string {
	return "Модуль для генерации цитат пользователей"
}

func (w *QuoteModule) Commands() []core.Command {
	return []core.Command{}
}

type quoteCommand struct{}

func (c *quoteCommand) Info() *core.CommandInfo {
	return &core.CommandInfo{
		Name: "Quote",
		Desc: "Генерирует цитату",
	}
}

func (c *quoteCommand) Usage() *core.CommandUsage {
	return &core.CommandUsage{
		Desc:   "Генерирует цитату из сообщения пользователя",
		Params: []core.CommandUsageParam{
            {Name: "dark", Desc: "Меняет цвет фона цитаты на чёрный", Optional: true},
        },
	}
}

func (c *quoteCommand) Run(msg events.MessageNewObject, args []string, bot *core.Bot) {
}
