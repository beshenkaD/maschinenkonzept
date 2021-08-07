package main

import (
	"os"

	_ "github.com/beshenkaD/maschinenkonzept/admin"
	_ "github.com/beshenkaD/maschinenkonzept/config"
	"github.com/beshenkaD/maschinenkonzept/core"
	_ "github.com/beshenkaD/maschinenkonzept/me"
	_ "github.com/beshenkaD/maschinenkonzept/ping"
)

func main() {
	core.AddLanguage("ru", &core.Strings{
		CommandDisabled:       "команда отключена в этом чате",
		InvalidCommand:        "неправильная команда",
		Error:                 "Ошибка: ",
		HelpDescription:       "Описание: %s",
		HelpUsage:             "Пример: %s%s %s",
		HelpAvailableCommands: "Доступные команды: %v",
		HelpAboutCommand:      "Введите: '%shelp <команда>' чтобы увидеть помощь для конкретной команды.",
	})

	core.AddLanguage("uk", &core.Strings{
		CommandDisabled:       "команда відключена в цьому чаті",
		InvalidCommand:        "неправильна команда",
		Error:                 "Помилка: ",
		HelpDescription:       "Опис: %s",
		HelpUsage:             "Приклад: %s%s %s",
		HelpAvailableCommands: "Доступні команди: %v",
		HelpAboutCommand:      "Введіть: '%shelp <команда>' щоб побачити допомога для конкретної команди",
	})

	bot := core.New(os.Getenv("VK_TOKEN"), "/home/beshenka/hueta", true)

	bot.Run()
}
