package vkutil

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"math/rand"
	"time"
)

/*
Простая обёртка для отправки сообщений
*/
func SendMessage(session *api.VK, message string, peerID int, disableMentions bool) (int, error) {
	rand.Seed(time.Now().UnixNano())

	b := params.NewMessagesSendBuilder()
	b.Message(message)
	b.RandomID(rand.Int())
	b.PeerID(peerID)
	b.DisableMentions(disableMentions)

	i, err := session.MessagesSend(b.Params)
	return i, err
}
