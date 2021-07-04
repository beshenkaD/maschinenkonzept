// Базовый модуль бота

package base

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/beshenkaD/maschinenkonzept/core"
	"github.com/beshenkaD/maschinenkonzept/vkutil"
	"runtime"
	"syscall"
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
	vkutil.SendMessage(bot.Session, "Привет! 👋", msg.Message.PeerID, true)
}

func (w *BaseModule) OnKickUser(bot *core.Bot, msg events.MessageNewObject) {
	vkutil.SendMessage(bot.Session, "Пока 👋", msg.Message.PeerID, true)
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
	vkutil.SendMessage(bot.Session, "pong", msg.Message.PeerID, true)
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
⚙ OS: %s
⚙ Uptime: %s
⚙ Сообщений обработано: %d
⚙ Потребление памяти (alloc): %v MiB
`
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	v := m.Alloc / 1024 / 1024
	u := time.Since(bot.StartTime)
	os := runtime.GOOS

	s = fmt.Sprintf(s, bot.SelfName, os, u, bot.Processed, v)

	var r syscall.Rusage
	err := syscall.Getrusage(syscall.RUSAGE_SELF, &r)

	if err == nil {
		s += fmt.Sprintf("⚙ Потребление памяти (rusage): %v MiB", r.Maxrss/1024)
	}

	vkutil.SendMessage(bot.Session, s, msg.Message.PeerID, true)
}

func (c *statCommand) Usage() *core.CommandUsage {
	return &core.CommandUsage{
		Desc:   "Выводит статистику бота",
		Params: []core.CommandUsageParam{},
	}
}
