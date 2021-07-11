package core

import (
	"fmt"
	"strings"

	"github.com/beshenkaD/maschinenkonzept/vkutil"
)

type ConfigModule struct{}

func (w *ConfigModule) Name() string {
	return "Конфигурация"
}

func (w *ConfigModule) Commands() []Command {
	return []Command{
		&getConfigCommand{},
		&setupCommand{},
	}
}

func (w *ConfigModule) Description() string {
	return "Управление конфигурацией беседы"
}

func (w *ConfigModule) OnInviteBot(chat *Chat, msg vkMessage) {
	vkutil.SendMessage(chat.Bot.Session, "Привет! Настрой меня командой /setup", chat.ID, true)
}

type getConfigCommand struct{}

func (c *getConfigCommand) Info() *CommandInfo {
	return &CommandInfo{
		Name:    "GetConfig",
		Desc:    "fdsfsdf",
		ForConf: true,
		ForPm:   true,
	}
}

func (c *getConfigCommand) Usage() *CommandUsage {
	return &CommandUsage{}
}

func (c *getConfigCommand) Run(msg vkMessage, args []string, chat *Chat) string {
	config := chat.Config

	s := fmt.Sprintf(`Настройка командой /setup выполнена: %t
Базовые настройки:
-- Игнорировать неправильные команды: %t
-- Алиасы: TODO
-- Префикс для команд: %s

Настройки модулей:
-- Отключенные модули: TODO
-- Отключенные команды: TODO
`, config.SetupDone, config.Basic.IgnoreInvalidCommands, config.Basic.CommandPrefix)

	return s
}

type setupCommand struct{}

func (c *setupCommand) disableModule(chat *Chat, module string) {
	for _, v := range chat.Modules {
		if strings.ToLower(v.Name()) == module {
			cmds := v.Commands()
			for _, c := range cmds {
				str := strings.ToLower(c.Info().Name)
				d := &chat.Config.Modules.CommandDisabled

				if len(*d) <= 0 {
					*d = make(map[string]bool)
				}
				chat.Config.Modules.CommandDisabled[str] = true
			}
		}
	}
	chat.Config.Modules.Disabled[module] = true
}

func (c *setupCommand) Info() *CommandInfo {
	return &CommandInfo{
		Name:    "Setup",
		Desc:    "Настраивает бота для текущего чата",
		ForConf: true,
		ForPm:   false,
	}
}

func (c *setupCommand) Usage() *CommandUsage {
	return &CommandUsage{
		Desc: "Эта команда поможет провести настройку бота для текущего чата",
		Params: []CommandUsageParam{
			{Name: "Override", Desc: "Перезаписывает текущие настройки", Optional: true},
			{Name: "Префикс", Desc: "Префикс с помощью которого вы сможете обращаться к боту", Optional: true},
			{Name: "Игнорировать неправильные команды", Desc: "[true] или [false]", Optional: true},
		},
	}
}

func (c *setupCommand) Run(msg vkMessage, args []string, chat *Chat) string {
	if len(args) < 1 {
		chat.Config = *DefaultConfig()
		chat.Config.SetupDone = true
		return "Вы не передали никаких аргументов. Использую конфигурацию по-умолчанию"
	}

	if chat.Config.SetupDone {
		if strings.ToLower(args[0]) != "override" {
			return "Беседа уже была сконфигурирована! Используйте override чтобы перезаписать настройки"
		}
		args = args[1:]
		chat.Config = *DefaultConfig()
	}

	if len(args) >= 1 {
		chat.Config.Basic.CommandPrefix = args[0]
	}

	if len(args) >= 2 {
		s := args[1]
		switch strings.ToLower(s) {
		case "true":
			chat.Config.Basic.IgnoreInvalidCommands = true
		case "false":
			chat.Config.Basic.IgnoreInvalidCommands = false
		default:
			return fmt.Sprintf("Неверный аргумент: %s", s)
		}
	}

	chat.Config.Basic.Aliases = make(map[string]string)
	chat.Config.Modules.CommandDisabled = make(map[string]bool)
	chat.Config.Modules.Disabled = make(map[string]bool)

	chat.Config.SetupDone = true

	return "Вы успешно настроили бота"
}
