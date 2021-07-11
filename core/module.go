package core

import (
	"github.com/SevereCloud/vksdk/v2/events"
)

type vkMessage events.MessageNewObject

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
	OnMessage(*Chat, vkMessage)
}

// Хук для добавления пользователя в беседу
type ModuleOnInviteUser interface {
	Module
	OnInviteUser(*Chat, vkMessage)
}

// Хук для кика пользователя из беседы
type ModuleOnKickUser interface {
	Module
	OnKickUser(*Chat, vkMessage)
}

// Хук для закрепления сообщения
type ModuleOnPinMessage interface {
	Module
	OnPinMessage(*Chat, vkMessage)
}

// Хук для открепления сообщения
type ModuleOnUnpinMessage interface {
	Module
	OnUnpinMessage(*Chat, vkMessage)
}

// Хук для вступления по ссылке
type ModuleOnInviteByLink interface {
	Module
	OnInviteByLink(*Chat, vkMessage)
}

// Хук для создания чата
type ModuleOnChatCreate interface {
	Module
	OnChatCreate(*Chat, vkMessage)
}

// Хук для приглашения бота (этого)
type ModuleOnInviteBot interface {
	Module
	OnInviteBot(*Chat, vkMessage)
}

// Хук выполняется каждую секунду
type ModuleOnTick interface {
	Module
	OnTick(*Chat)
}

type moduleHooks struct {
	OnMessage      []ModuleOnMessage
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
	Run(vkMessage, []string, *Chat) string
	Usage() *CommandUsage
	Info() *CommandInfo
}

func (c *Chat) RegisterModule(m Module) {
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
