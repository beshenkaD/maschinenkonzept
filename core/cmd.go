package core

// Parsed user input (used for commands)
type CommandInput struct {
	Message     string   // Full string without prefix
	MessageData *Message // Extra data about message
	Command     string   // First argument passed to the bot
	Args        []string // Arguments
	Chat        int      // Chat where the command was called (ID)
	User        *User    // User who sent the message
}

type HookInput struct {
	Raw         string
	MessageData *Message
	Chat        int
	User        *User
}

type Message struct {
	Text         string
	ActionType   string // https://vk.com/dev/objects/message look at action object
	MemberId     int
	ActionText   string
	FwdMessages  []Message
	ReplyMessage *Message
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

type Command struct {
	Name        string
	Cmd         string
	Func        CmdFunc
	Description string
	Params      []HelpParam
}

type Hook struct {
	Name        string
	ActionType  string
	Func        HookFunc
	Description string
}

type Tick struct {
	Name        string
	Func        TickFunc
	Description string
}

type CmdFunc func(in *CommandInput) (string, error)
type HookFunc func(in *HookInput) (string, error)
type TickFunc func() string

var (
	commands = make(map[string]*Command)
	ticks    = make(map[string]*Tick)
	hooks    = make(map[string]*Hook)
)

func RegisterCommand(name, trigger, description string, params []HelpParam, cmdFunc CmdFunc) {
	commands[name] = &Command{
		Name:        name,
		Cmd:         trigger,
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
		if !ignoreInvalid[i.Chat] {
			b.SendMessage(i.Chat, "Invalid command")
		}
		return
	}

	message, err := cmd.Func(i)
	if err != nil {
		b.ErrorHandler(i.Chat, err)
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

	if IsHookDisabled(hook.Name, i.Chat) {
		return
	}

	message, err := hook.Func(i)

	if err != nil {
		b.ErrorHandler(i.Chat, err)
		return
	}

	if message != "" {
		b.SendMessage(i.Chat, message)
	}
}

func (b *Bot) handleTick(name string, chat int) {
	if IsTickDisabled(name, chat) {
		return
	}

	tick, ok := ticks[name]

	if !ok {
		return
	}

	message := tick.Func()

	if message != "" {
		b.SendMessage(chat, message)
	}
}
