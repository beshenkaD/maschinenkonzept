package core

import (
	"strings"

	"github.com/SevereCloud/vksdk/v2/events"
)

// Модуль отслеживает все входящие запросы в зависимости от того какие интерфейсы он реализует
type Module interface {
	Name() string
	Commands() []Command
	Description() string
}

// Вопрос. Нахуя так много интерфейсов? Ответ: чтобы можно было реализовывать их частично
// Хук для нового сообщения
type ModuleOnMessage interface {
	Module
	OnMessage(*Bot, events.MessageNewObject)
}

// Хук для команды (команды отделяются от обычных сообщений)
type ModuleOnCommand interface {
	Module
	OnCommand(*Bot, events.MessageNewObject)
}

// Хук для добавления пользователя в беседу
type ModuleOnInviteUser interface {
	Module
	OnInviteUser(*Bot, events.MessageNewObject)
}

// Хук для кика пользователя из беседы
type ModuleOnKickUser interface {
	Module
	OnKickUser(*Bot, events.MessageNewObject)
}

// Хук для закрепления сообщения
type ModuleOnPinMessage interface {
	Module
	OnPinMessage(*Bot, events.MessageNewObject)
}

// Хук для открепления сообщения
type ModuleOnUnpinMessage interface {
	Module
	OnUnpinMessage(*Bot, events.MessageNewObject)
}

// Хук для вступления по ссылке
type ModuleOnInviteByLink interface {
	Module
	OnInviteByLink(*Bot, events.MessageNewObject)
}

// Хук для создания чата
type ModuleOnChatCreate interface {
	Module
	OnChatCreate(*Bot, events.MessageNewObject)
}

// Хук для приглашения бота (этого)
type ModuleOnInviteBot interface {
	Module
	OnInviteBot(*Bot, events.MessageNewObject)
}

// Хук выполняется каждую секунду
type ModuleOnTick interface {
	Module
	OnTick(*Bot)
}

type moduleHooks struct {
	OnMessage      []ModuleOnMessage
	OnCommand      []ModuleOnCommand
	OnInviteUser   []ModuleOnInviteUser
	OnKickUser     []ModuleOnKickUser
	OnPinMessage   []ModuleOnPinMessage
	OnUnpinMessage []ModuleOnUnpinMessage
	OnInviteByLink []ModuleOnInviteByLink
	OnChatCreate   []ModuleOnChatCreate
	OnInviteBot    []ModuleOnInviteBot
	OnTick         []ModuleOnTick
}

// -------------------------------- //

// Описывает каждый параметр команды
type CommandUsageParam struct {
	Name     string
	Desc     string
	Optional bool
}

// Хелп для команды
type CommandUsage struct {
	Desc   string
	Params []CommandUsageParam
}

// Информация о команде
type CommandInfo struct {
	Name string
	Desc string
}

// Команда это любая команда адресованная боту
type Command interface {
	Run(events.MessageNewObject, []string, *Bot)
	Usage() *CommandUsage
	Info() *CommandInfo
	ForPm() bool
	ForConf() bool
}

func (b *Bot) RegisterModule(m Module) {
	cmds := m.Commands()
	for _, c := range cmds {
		b.addCommand(c, m)
	}

	if h, ok := m.(ModuleOnCommand); ok {
		b.hooks.OnCommand = append(b.hooks.OnCommand, h)
	}

	if h, ok := m.(ModuleOnMessage); ok {
		b.hooks.OnMessage = append(b.hooks.OnMessage, h)
	}

	if h, ok := m.(ModuleOnInviteUser); ok {
		b.hooks.OnInviteUser = append(b.hooks.OnInviteUser, h)
	}

	if h, ok := m.(ModuleOnKickUser); ok {
		b.hooks.OnKickUser = append(b.hooks.OnKickUser, h)
	}

	if h, ok := m.(ModuleOnPinMessage); ok {
		b.hooks.OnPinMessage = append(b.hooks.OnPinMessage, h)
	}

	if h, ok := m.(ModuleOnUnpinMessage); ok {
		b.hooks.OnUnpinMessage = append(b.hooks.OnUnpinMessage, h)
	}

	if h, ok := m.(ModuleOnInviteByLink); ok {
		b.hooks.OnInviteByLink = append(b.hooks.OnInviteByLink, h)
	}

	if h, ok := m.(ModuleOnChatCreate); ok {
		b.hooks.OnChatCreate = append(b.hooks.OnChatCreate, h)
	}

	if h, ok := m.(ModuleOnInviteBot); ok {
		b.hooks.OnInviteBot = append(b.hooks.OnInviteBot, h)
	}

	if h, ok := m.(ModuleOnTick); ok {
		b.hooks.OnTick = append(b.hooks.OnTick, h)
	}
}

func (b *Bot) addCommand(c Command, m Module) {
	name := strings.ToLower(c.Info().Name)
	b.commands[name] = c
}
