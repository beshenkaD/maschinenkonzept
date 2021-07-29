package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/beshenkaD/maschinenkonzept/vkutil"
)

type Chat struct {
	ID         int
	ConfigLock sync.Mutex
	Config     Config
	hooks      moduleHooks
	Modules    []Module
	commands   map[commandID]Command
	Bot        *Bot
}

func NewChat(bot *Bot, ID int) *Chat {
	return &Chat{
		ID:       ID,
		commands: make(map[commandID]Command),
		Config:   *NewConfig(),
		Bot:      bot,
	}
}

func (ch *Chat) SendMessage(msg string) {
	vkutil.SendMessage(ch.Bot.Session, msg, ch.ID, false)
}

func (ch *Chat) RemoveUser() {

}

/*
   Использовать json файлы оказалось проще чем нормальную бд
   Причины:
      1. В дальнейшем в конфиг будут добавляться новые поля. В sql базах данных это ебаный геморрой
      2. Адекватной sql базы данных (sqlite) под гошку без ебли нема

   Я надеюсь это временное решение. В будущем можно приспособить сюда mysql или mongo
*/

// Сохраняет конфигурацию беседы на диск
func (ch *Chat) WriteConfig() {
	bs, err := json.Marshal(ch.Config)

	if err != nil {
		return
	}

	f := path.Join(ch.Bot.ConfigsPath, fmt.Sprintf("%d.json", ch.ID))
	if err := os.WriteFile(f, bs, 0664); err != nil {
		return
	}
}

// Загружает конфигурацию беседы с диска
func (ch *Chat) LoadConfig() error {
	f := path.Join(ch.Bot.ConfigsPath, fmt.Sprintf("%d.json", ch.ID))
	content, err := os.ReadFile(f)

	if err != nil {
		return err
	}

	var config Config
	err = json.Unmarshal(content, &config)

	if err != nil {
		return err
	}

	ch.Config = config

	return nil
}

func (ch *Chat) addCommand(c Command, m Module) {
	name := strings.ToLower(c.Info().Name)
	ch.commands[commandID(name)] = c
}

func (ch *Chat) ShouldRunHooks(m Module) bool {
	if ch.ID < 2000000000 {
		return false
	}

	if _, off := ch.Config.Modules.Disabled[moduleID(strings.ToLower(m.Name()))]; off {
		return false
	}

	return true
}
