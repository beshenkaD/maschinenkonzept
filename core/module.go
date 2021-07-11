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

// ""
type ModuleOnMessage interface {
	Module
	OnMessage(*Chat, vkMessage)
}

// chat_photo_update
type ModuleOnPhotoUpdate interface {
	Module
	OnPhotoUpdate(*Chat, vkMessage)
}

// chat_photo_remove
type ModuleOnPhotoRemove interface {
	Module
	OnPhotoRemove(*Chat, vkMessage)
}

// chat_create
type ModuleOnChatCreate interface {
	Module
	OnChatCreate(*Chat, vkMessage)
}

// chat_title_update
type ModuleOnTitleUpdate interface {
	Module
	OnTitleUpdate(*Chat, vkMessage)
}

// chat_invite_user
type ModuleOnInviteUser interface {
	Module
	OnInviteUser(*Chat, vkMessage)
}

// chat_kick_user
type ModuleOnKickUser interface {
	Module
	OnKickUser(*Chat, vkMessage)
}

// chat_pin_message
type ModuleOnPinMessage interface {
	Module
	OnPinMessage(*Chat, vkMessage)
}

// chat_unpin_message
type ModuleOnUnpinMessage interface {
	Module
	OnUnpinMessage(*Chat, vkMessage)
}

// chat_invite_user_by_link
type ModuleOnInviteByLink interface {
	Module
	OnInviteByLink(*Chat, vkMessage)
}

type ModuleOnTick interface {
	Module
	OnTick(*Chat)
}

type moduleHooks struct {
	OnMessage      []ModuleOnMessage
	OnPhotoUpdate  []ModuleOnPhotoUpdate
	OnPhotoRemove  []ModuleOnPhotoRemove
	OnChatCreate   []ModuleOnChatCreate
	OnTitleUpdate  []ModuleOnTitleUpdate
	OnInviteUser   []ModuleOnInviteUser
	OnKickUser     []ModuleOnKickUser
	OnPinMessage   []ModuleOnPinMessage
	OnUnpinMessage []ModuleOnUnpinMessage
	OnInviteByLink []ModuleOnInviteByLink
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

	if h, ok := m.(ModuleOnPhotoUpdate); ok {
		c.hooks.OnPhotoUpdate = append(c.hooks.OnPhotoUpdate, h)
	}

	if h, ok := m.(ModuleOnPhotoRemove); ok {
		c.hooks.OnPhotoRemove = append(c.hooks.OnPhotoRemove, h)
	}

	if h, ok := m.(ModuleOnTitleUpdate); ok {
		c.hooks.OnTitleUpdate = append(c.hooks.OnTitleUpdate, h)
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

	if h, ok := m.(ModuleOnTick); ok {
		c.hooks.OnTick = append(c.hooks.OnTick, h)
	}
}
