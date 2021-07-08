package core

import (
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
	Name    string
	Desc    string
	ForPm   bool
	ForConf bool
}

// Команда это любая команда адресованная боту
type Command interface {
	Run(events.MessageNewObject, []string, *Bot)
	Usage() *CommandUsage
	Info() *CommandInfo
}

func (c *Conversation) RegisterModule(m Module) {
	if h, ok := m.(ModuleOnCommand); ok {
		c.hooks.OnCommand = append(c.hooks.OnCommand, h)
	}

	if h, ok := m.(ModuleOnMessage); ok {
		c.hooks.OnMessage = append(c.hooks.OnMessage, h)
	}

	if h, ok := m.(ModuleOnInviteUser); ok {
		c.hooks.OnInviteUser = append(c.hooks.OnInviteUser, h)
	}

	if h, ok := m.(ModuleOnKickUser); ok {
		c.hooks.OnKickUser = append(c.hooks.OnKickUser, h)
	}

	if h, ok := m.(ModuleOnPinMessage); ok {
		c.hooks.OnPinMessage = append(c.hooks.OnPinMessage, h)
	}

	if h, ok := m.(ModuleOnUnpinMessage); ok {
		c.hooks.OnUnpinMessage = append(c.hooks.OnUnpinMessage, h)
	}

	if h, ok := m.(ModuleOnInviteByLink); ok {
		c.hooks.OnInviteByLink = append(c.hooks.OnInviteByLink, h)
	}

	if h, ok := m.(ModuleOnChatCreate); ok {
		c.hooks.OnChatCreate = append(c.hooks.OnChatCreate, h)
	}

	if h, ok := m.(ModuleOnInviteBot); ok {
		c.hooks.OnInviteBot = append(c.hooks.OnInviteBot, h)
	}

	if h, ok := m.(ModuleOnTick); ok {
		c.hooks.OnTick = append(c.hooks.OnTick, h)
	}
}
