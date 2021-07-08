package core

import (
	"log"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api"
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
	var response api.MessagesGetConversationMembersResponse

	params := api.Params{
		"peer_id": ID,
	}

	err := bot.Session.RequestUnmarshal("messages.getConversationMembers", &response, params)
	if err != nil {
		log.Println(err)
	}

	var owner int
	for i, u := range response.Items {
		if u.IsOwner {
			owner = response.Profiles[i].ID
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
func (info *Conversation) ShouldRunModule(chatID int, m Module) bool {
	return true
}
