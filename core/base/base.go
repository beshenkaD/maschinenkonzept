// Базовый модуль бота

package base

import (
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/beshenkaD/maschinenkonzept/apiutil"
	"github.com/beshenkaD/maschinenkonzept/core"
	"runtime"
    "fmt"
    "time"
)

type BaseModule struct{}

func New() *BaseModule {
	return &BaseModule{}
}

func (w *BaseModule) Name() string {
	return "Основа"
}

func (w *BaseModule) OnInviteUser(bot *core.Bot, msg events.MessageNewObject) {
	apiutil.Send(bot.Session, "Привет! 👋", msg.Message.PeerID)
}

func (w *BaseModule) OnKickUser(bot *core.Bot, msg events.MessageNewObject) {
	apiutil.Send(bot.Session, "Пока 👋", msg.Message.PeerID)
}

func (w *BaseModule) Commands() []core.Command {
	return []core.Command{
		&pingCommand{},
        &statCommand{},
	}
}

func (w *BaseModule) Description() string {
	return "Базовый модуль для проверки работоспособности бота"
}

type pingCommand struct{}

func (c *pingCommand) Info() *core.CommandInfo {
	return &core.CommandInfo{
		Name: "Ping",
		Desc: "Проверяет работоспособность бота и позволяет поиграть с ним в пинг-понг⚾",
	}
}

func (c *pingCommand) Run(msg events.MessageNewObject, argc int, argv []string, bot *core.Bot) {
	apiutil.Send(bot.Session, "pong", msg.Message.PeerID)
}

func (c *pingCommand) Usage() *core.CommandUsage {
	return &core.CommandUsage{
		Desc:   "Проверяет работоспособность бота",
		Params: []core.CommandUsageParam{},
	}
}

type statCommand struct{}

func (c *statCommand) Info() *core.CommandInfo {
	return &core.CommandInfo{
		Name: "Stat",
		Desc: "Выводит статистику бота 🚀",
	}
}

func (c *statCommand) Run(msg events.MessageNewObject, argc int, argv []string, bot *core.Bot) {
	s := `⚙ Запущен как: %s
⚙ Uptime: %s
⚙ Сообщений обработано: %d
⚙ Потребление памяти (alloc): %v MiB
`
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    v := m.Alloc / 1024 / 1024
    u := time.Since(bot.StartTime)

    s = fmt.Sprintf(s, bot.SelfName, u, bot.Processed, v)

    apiutil.Send(bot.Session, s, msg.Message.PeerID)
}

func (c *statCommand) Usage() *core.CommandUsage {
	return &core.CommandUsage{
		Desc: "Выводит статистику бота",
		Params: []core.CommandUsageParam{},
	}
}
