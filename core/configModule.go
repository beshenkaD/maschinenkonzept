package core

import (
	"fmt"
	"strings"

	"github.com/beshenkaD/maschinenkonzept/vkutil"
)

type ConfigModule struct{}

func (w *ConfigModule) Name() string {
	return "config"
}

func (w *ConfigModule) Commands() []Command {
	return []Command{
		&getConfigCommand{},
		&setupCommand{},
		&setConfigCommand{},
	}
}

func (w *ConfigModule) Description() string {
	return "Управление конфигурацией беседы"
}

func (w *ConfigModule) OnInviteUser(chat *Chat, msg vkMessage) {
	if msg.Message.Action.MemberID == -chat.Bot.SelfID {
		vkutil.SendMessage(chat.Bot.Session, "Привет! Настрой меня командой /setup", chat.ID, true)
	}
}

func disableModule(chat *Chat, module string) {
	for _, v := range chat.Modules {
		if strings.ToLower(v.Name()) == module {
			cmds := v.Commands()
			for _, c := range cmds {
				str := strings.ToLower(c.Info().Name)
				d := &chat.Config.Modules.DisabledCommands

				if len(*d) <= 0 {
					*d = make(map[commandID]bool)
				}
				chat.Config.Modules.DisabledCommands[commandID(str)] = true
			}
		}
	}
	chat.Config.Modules.Disabled[moduleID(module)] = true
}

type setConfigCommand struct{}

func (c *setConfigCommand) Info() *CommandInfo {
	return &CommandInfo{
		Name:    "SetConfig",
		Desc:    "Устанавливает значение в конфигурации и сохраняет её",
		ForConf: true,
		ForPm:   true,
	}
}

func (c *setConfigCommand) Usage() *CommandUsage {
	return &CommandUsage{
		Desc: "Устанавливает значние в формате `коллекция.параметр`. Чтобы удалить значение установите его в пустоту",
		Params: []CommandUsageParam{
			{Name: "[параметр] [значение]", Desc: "", Optional: true},
			{Name: "[словарь] [ключ] [значение]", Desc: "", Optional: true},
		},
	}
}

func (c *setConfigCommand) Run(msg vkMessage, args []string, chat *Chat) string {
	if len(args) < 1 {
		return "Вы не передали никаких параметров!"
	}

	if len(args) < 2 {
		return "Вы не передали никаких значений!"
	}

	s, ok := chat.Config.Set(chat, args, msg.Message.Text)

	if !ok {
		return s
	}

	go chat.WriteConfig()

	return s
}

type getConfigCommand struct{}

func (c *getConfigCommand) Info() *CommandInfo {
	return &CommandInfo{
		Name:    "GetConfig",
		Desc:    "Возвращает текущую конфигурацию бота",
		ForConf: true,
		ForPm:   true,
	}
}

func (c *getConfigCommand) Usage() *CommandUsage {
	return &CommandUsage{}
}

func (c *getConfigCommand) Run(msg vkMessage, args []string, chat *Chat) string {
	config := chat.Config

	m := ""
	cm := ""
	al := ""

	for module := range chat.Config.Modules.Disabled {
		m += string(module) + "\n"
	}

	for cmd := range chat.Config.Modules.DisabledCommands {
		cm += string(cmd) + "\n"
	}

	for key, value := range chat.Config.Basic.Aliases {
		al += key + " -> " + value + "\n"
	}

	s := fmt.Sprintf(`Настройка командой setup выполнена: %s
Базовые настройки:
-- Игнорировать неправильные команды: %s
-- Алиасы: %s
-- Префикс для команд: %s

Настройки модулей:
-- Отключенные модули: %s
-- Отключенные команды: %s
`, boolToRus(config.SetupDone), boolToRus(config.Basic.IgnoreInvalidCommands), al, config.Basic.CommandPrefix, m, cm)

	return s
}

type setupCommand struct{}

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
		chat.Config = *NewConfig()
		chat.Config.SetupDone = true

		return "Вы не передали никаких аргументов. Использую конфигурацию по-умолчанию"
	}

	if chat.Config.SetupDone {
		if strings.ToLower(args[0]) != "override" {
			return "Беседа уже была сконфигурирована! Используйте override чтобы перезаписать настройки"
		}
		args = args[1:]
		chat.Config = *NewConfig()
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
	chat.Config.Modules.DisabledCommands = make(map[commandID]bool)
	chat.Config.Modules.Disabled = make(map[moduleID]bool)

	chat.Config.SetupDone = true

	go chat.WriteConfig()

	return "Вы успешно настроили бота"
}
