package apiutil

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"math/rand"
	"time"
)

func Send(session *api.VK, message string, peer_id int) (int, error) {
	b := params.NewMessagesSendBuilder()
	b.Message(message)

	rand.Seed(time.Now().UnixNano())
	b.RandomID(rand.Int())

	b.PeerID(peer_id)

	i, err := session.MessagesSend(b.Params)
	return i, err
}
