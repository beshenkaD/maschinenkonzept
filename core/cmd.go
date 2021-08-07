package core

import (
	"errors"

	"github.com/SevereCloud/vksdk/v2/object"
)

// TODO: create more convinient structure for message
type Message object.MessagesMessage

type User struct {
	ID        int
	FirstName string
	LastName  string
	IsBot     bool
}

// Parsed user input (used for commands)
type CommandInput struct {
	Command string   // First argument passed to the bot
	Args    []string // Arguments
	Message *Message // Raw message
	Chat    *Chat    // Chat where the command was called
	User    *User    // User who sent the message
}

// User input for hooks
type HookInput struct {
	Message *Message // Raw message
	Chat    *Chat    // Chat where message was sent
	User    *User    // User who sent the message
}

type HelpParam struct {
	Name        string
	Description string
	Optional    bool
}

// A command is a message that is called by the user using a prefix and a trigger
// Example: /ping
// The command can have arguments. All arguments must be described in the Params array
type Command struct {
	Name        string
	Trigger     string
	Func        CmdFunc
	Description string
	Params      []HelpParam
}

// A hook is a passive command that is called by some event, or every time the bot receives a message
// List of available action types:
// - ActionNewMessage
// - ActionPhotoUpdate
// - ActionPhotoRemove
// - ActionChatCreate
// - ActionTitleUpdate
// - ActionInviteUser
// - ActionKickUser
// - ActionPinMessage
// - ActionUnpinMessage
// - ActionInviteByLink

type Hook struct {
	Name        string
	ActionType  actionType
	Func        HookFunc
	Description string
}

// Tick is a passive command that is executed every 5 seconds only in active chats
type Tick struct {
	Name        string
	Func        TickFunc
	Description string
}

type CmdFunc func(in *CommandInput) (string, error)
type HookFunc func(in *HookInput) (string, error)
type TickFunc func(chat *Chat) string

var (
	commands = make(map[string]*Command)
	hooks    = make(map[actionType][]*Hook)
	ticks    = make(map[string]*Tick)
)

func RegisterCommand(name, trigger, description string, params []HelpParam, cmdFunc CmdFunc) {
	commands[trigger] = &Command{
		Name:        name,
		Trigger:     trigger,
		Func:        cmdFunc,
		Description: description,
		Params:      params,
	}
}

func RegisterHook(name, description string, action actionType, hookFunc HookFunc) {
	hooks[action] = append(hooks[action], &Hook{
		Name:        name,
		ActionType:  action,
		Func:        hookFunc,
		Description: description,
	})
}

func RegisterTick(name, description string, periodicFunc TickFunc) {
	ticks[name] = &Tick{
		Name:        name,
		Func:        periodicFunc,
		Description: description,
	}
}

func (b *Bot) handleCommand(i *CommandInput) {
	cmd := commands[i.Command]

	if cmd == nil {
		if !i.Chat.IgnoreInvalid {
			b.sendError(i.Chat, errors.New("invalid command"))
		}
		return
	}

	message, err := cmd.Func(i)
	if err != nil {
		b.sendError(i.Chat, err)
		return
	}

	if message != "" {
		b.sendMessage(i.Chat, message)
	}
}

func (b *Bot) handleHook(i *HookInput) {
	hooks := hooks[parseAction(i.Message.Action.Type)]

	if len(hooks) == 0 {
		return
	}

	for _, hook := range hooks {
		if i.Chat.IsHookDisabled(hook.Name) {
			return
		}

		message, err := hook.Func(i)

		if err != nil {
			b.sendError(i.Chat, err)
			return
		}

		b.sendMessage(i.Chat, message)
	}
}

func (b *Bot) handleTick(name string, chat *Chat) {
	if chat.IsTickDisabled(name) {
		return
	}

	tick, ok := ticks[name]

	if !ok {
		return
	}

	message := tick.Func(chat)

	b.sendMessage(chat, message)
}
