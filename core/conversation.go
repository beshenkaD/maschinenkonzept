package core

import (
	"log"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api/params"
)

type Conversation struct {
	ID      int
	OwnerId int
	// TODO: Name string
	// TODO: Config BotConfig
	hooks    moduleHooks
	Prefix   byte
	Modules  []Module
	commands map[string]Command
	Bot      *Bot
}

// Возвращает объект беседы с дефолтной конфигурацией
func NewConversation(bot *Bot, ID int) *Conversation {
	b := params.NewMessagesGetConversationMembersBuilder()
	b.PeerID(ID)

	users, err := bot.Session.MessagesGetConversationMembers(b.Params)

	if err != nil {
		log.Println(err.Error())
	}

	var owner int
	for _, u := range users.Items {
		if u.IsOwner {
			owner = u.MemberID
		}
	}

	return &Conversation{
		ID:       ID,
		OwnerId:  owner,
		commands: make(map[string]Command),
		Prefix:   '/',
		Bot:      bot,
	}
}

func (info *Conversation) addCommand(c Command, m Module) {
	name := strings.ToLower(c.Info().Name)
	info.commands[name] = c
}

// TODO
func (info *Conversation) ShouldRunHooks(chatID int, m Module) bool {
	if chatID < 2000000000 {
		return false
	}

	return true
}
