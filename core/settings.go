package core

var (
	ignoreInvalid    = make(map[int]bool)
	prefix           = make(map[int]byte)
	language         = make(map[int]string)
	disabledCommands = make(map[int][]string)
	disabledTicks    = make(map[int][]string)
	disabledHooks    = make(map[int][]string)
)

const (
	DefaultLanguage = "en"
)

func IgnoreInvalid(b bool, chat int) {
	ignoreInvalid[chat] = b
}

func Prefix(p byte, chat int) {
	prefix[chat] = p
}

func DisableCommand(cmd string, chat int) {
	disabledCommands[chat] = append(disabledCommands[chat], cmd)
}

func DisableHook(hook string, chat int) {
	disabledHooks[chat] = append(disabledHooks[chat], hook)
}

func DisableTick(hook string, chat int) {
	disabledTicks[chat] = append(disabledTicks[chat], hook)
}

func EnableCommand(cmd string, chat int) {
	s := disabledCommands[chat]

	for i, c := range s {
		if c == cmd {
			s[i] = s[len(s)-1]
			s = s[:len(s)-1]
		}
	}

	disabledCommands[chat] = s
}

func EnableHook(hook string, chat int) {
	s := disabledHooks[chat]

	for i, c := range s {
		if c == hook {
			s[i] = s[len(s)-1]
			s = s[:len(s)-1]
		}
	}

	disabledHooks[chat] = s
}

func EnableTick(tick string, chat int) {
	s := disabledTicks[chat]

	for i, c := range s {
		if c == tick {
			s[i] = s[len(s)-1]
			s = s[:len(s)-1]
		}
	}

	disabledTicks[chat] = s
}

func isDisabled(t string, arg string, chat int) bool {
	var m map[int][]string

	switch t {
	case "command":
		m = disabledCommands
	case "hook":
		m = disabledHooks
	case "tick":
		m = disabledTicks
	}

	if g, ok := m[chat]; ok {
		for _, e := range g {
			if arg == e {
				return true
			}
		}
	}

	return false
}

func IsCommandDisabled(command string, chat int) bool {
	return isDisabled("command", command, chat)
}

func IsHookDisabled(hook string, chat int) bool {
	return isDisabled("hook", hook, chat)
}

func IsTickDisabled(tick string, chat int) bool {
	return isDisabled("tick", tick, chat)
}
