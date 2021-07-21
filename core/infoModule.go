package core

type InfoModule struct{}

func (w *InfoModule) Name() string {
	return "Info"
}

func (w *InfoModule) Commands() []Command {
	return []Command{
		&helpCommand{},
	}
}

func (w *InfoModule) Description() string {
	return "информация"
}

type helpCommand struct{}

func (c *helpCommand) Info() *CommandInfo {
	return &CommandInfo{
		Name:    "help",
		Desc:    "help",
		ForPm:   true,
		ForConf: true,
	}
}

func (c *helpCommand) Usage() *CommandUsage {
	return &CommandUsage{}
}

func (c *helpCommand) Run(msg vkMessage, args []string, chat *Chat) string {
	str := ""

	for _, m := range chat.Modules {
		str += "module " + m.Name() + "\n"
	}

	for _, c := range chat.commands {
		str += "command " + c.Info().Name + "\n"
	}

	return str
}
