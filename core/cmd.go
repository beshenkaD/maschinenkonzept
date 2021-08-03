package core

import "errors"

// Parsed user input (used for commands)
type CommandInput struct {
	Message     string   // Full string without prefix
	MessageData *Message // Extra data about message
	Command     string   // First argument passed to the bot
	Args        []string // Arguments
	Chat        *Chat    // Chat where the command was called
	User        *User    // User who sent the message
}

// User input for hooks
type HookInput struct {
	Raw         string   // Raw user message
	MessageData *Message // Extra data about message
	Chat        *Chat    // Chat where message was sent
	User        *User    // User who sent the message
}

// TODO: handle telegram actions
type Message struct {
	Text         string     // Actual text
	ActionType   string     // https://vk.com/dev/objects/message look at action object
	MemberId     int        // Vk specific thing
	ActionText   string     // Vk specific thing
	FwdMessages  []*Message // Forward messages (if any)
	ReplyMessage *Message   // Reply message
	IsPrivate    bool
}

type User struct {
	ID        int
	FirstName string
	LastName  string
	IsBot     bool
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
// The list of available action types can be found here: https://vk.com/dev/objects/message (action object)
// If you want the hook to be called on every message you must add a custom action type
// TODO: handle telegram events
type Hook struct {
	Name        string
	ActionType  string
	Func        HookFunc
	Description string
}

// Tick is a passive command that is executed every 5 seconds
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
	hooks    = make(map[string]*Hook)
	ticks    = make(map[string]*Tick)
)

func RegisterCommand(name, trigger, description string, params []HelpParam, cmdFunc CmdFunc) {
	commands[name] = &Command{
		Name:        name,
		Trigger:     trigger,
		Func:        cmdFunc,
		Description: description,
		Params:      params,
	}
}

func RegisterHook(name, action, description string, hookFunc HookFunc) {
	hooks[action] = &Hook{
		Name:        name,
		ActionType:  action,
		Func:        hookFunc,
		Description: description,
	}
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
			b.ErrorHandler(i.Chat.ID, errors.New("invalid command"))
		}
		return
	}

	message, err := cmd.Func(i)
	if err != nil {
		b.ErrorHandler(i.Chat.ID, err)
		return
	}

	if message != "" {
		b.SendMessage(i.Chat, message)
	}
}

func (b *Bot) handleHook(i *HookInput) {
	hook := hooks[i.MessageData.ActionType]

	if hook == nil {
		return
	}

	if i.Chat.IsHookDisabled(hook.Name) {
		return
	}

	message, err := hook.Func(i)

	if err != nil {
		b.ErrorHandler(i.Chat.ID, err)
		return
	}

	if message != "" {
		b.SendMessage(i.Chat, message)
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

	if message != "" {
		b.SendMessage(chat, message)
	}
}
