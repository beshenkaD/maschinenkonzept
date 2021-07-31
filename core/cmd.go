package core

// Parsed user input
type Input struct {
	Message     string // Full string without prefix
	MessageData *Message
	Command     string
	Args        []string
	Chat        int
	User        *User
}

type Message struct {
	Text       string
	ActionType string
	IsPrivate  bool
}

type User struct {
	ID    int
	Name  string
	IsBot bool
}

type CmdFunc func(in *Input) (string, error)
type HookFunc func(in *Input) (string, error)

// type PeriodicFunc func(duration) (string, error)

type Command struct {
	Name        string
	Cmd         string
	Func        CmdFunc
	Description string
	// Args
}

type Hook struct {
	Name        string
	ActionType  string // maybe int and use iota
	Func        HookFunc
	Description string
}

type PeriodicCommand struct {
	// Duration
	// Func
	Description string
}

var (
	commands = make(map[string]*Command)
	hooks    = make(map[string]*Hook)
)

func RegisterCommand(name, trigger, description string, cmdFunc CmdFunc) {
	commands[trigger] = &Command{
		Name:        name,
		Cmd:         trigger,
		Func:        cmdFunc,
		Description: description,
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

func RegisterPeriodicCommand() {}

func (b *Bot) DisableCommand(cmd string, chat int) {
	b.DisabledCommands[chat] = append(b.DisabledCommands[chat], cmd)
}

func (b *Bot) DisableHook(hook string, chat int) {
	b.DisabledHooks[chat] = append(b.DisabledHooks[chat], hook)
}

func (b *Bot) EnableCommand(cmd string, chat int) {
	s := b.DisabledCommands[chat]

	for i, c := range s {
		if c == cmd {
			s[i] = s[len(s)-1]
			s = s[:len(s)-1]
		}
	}

	b.DisabledCommands[chat] = s
}

func (b *Bot) EnableHook(hook string, chat int) {

}

func (b *Bot) IsCommandDisabled(command string, chat int) bool {
	if cmds, ok := b.DisabledCommands[chat]; ok {
		for _, cmd := range cmds {
			if command == cmd {
				return true
			}
		}
	}

	return false
}

func (b *Bot) handleCmd(i *Input) {
	cmd := commands[i.Command]

	if cmd == nil {
		// Handle error
		return
	}

	message, err := cmd.Func(i)
	if err != nil {
		// check
	}

	if message != "" {
		b.SendMessage(i.Chat, message)
	}
}

func (b *Bot) handleHook(i *Input) {
	hook := hooks[i.MessageData.ActionType]

	if hook == nil {
		return
	}

	message, err := hook.Func(i)

	if err != nil {
		// check
	}

	if message != "" {
		b.SendMessage(i.Chat, message)
	}
}
