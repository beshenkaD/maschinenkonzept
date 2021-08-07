package core

import "time"

// TODO: remove inactive chats from bot
type Chat struct {
	ID               int
	LastMessage      time.Time
	IgnoreInvalid    bool
	Prefix           string
	DisabledCommands map[string]bool
	DisabledHooks    map[string]bool
	DisabledTicks    map[string]bool
}

func newChat(ID int) *Chat {
	return &Chat{
		ID:               ID,
		LastMessage:      time.Now(),
		IgnoreInvalid:    false,
		Prefix:           "/",
		DisabledCommands: make(map[string]bool),
		DisabledHooks:    make(map[string]bool),
		DisabledTicks:    make(map[string]bool),
	}
}

func IsCommandExist(name string) bool {
	_, ok := commands[name]
	return ok
}

func IsHookExist(name string) bool {
	for _, hooks := range hooks {
		for _, hook := range hooks {
			if hook.Name == name {
				return true
			}
		}
	}

	return false
}

func IsTickExist(name string) bool {
	for _, tick := range ticks {
		if tick.Name == name {
			return true
		}
	}

	return false
}

func (c *Chat) IsCommandDisabled(name string) bool {
	t, ok := c.DisabledCommands[name]

	if ok {
		return t
	}

	return ok
}

func (c *Chat) IsHookDisabled(name string) bool {
	t, ok := c.DisabledHooks[name]

	if ok {
		return t
	}

	return ok
}

func (c *Chat) IsTickDisabled(name string) bool {
	t, ok := c.DisabledTicks[name]

	if ok {
		return t
	}

	return ok
}
