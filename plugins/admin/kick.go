package admin

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/beshenkaD/maschinenkonzept/apiutil"
)

func Kick(session *api.VK, message events.MessageNewObject) {
	apiutil.Send(session, "Пока-пока 😥", message.Message.PeerID)

	k := params.NewMessagesRemoveChatUserBuilder()
	k.ChatID(message.Message.PeerID - 2000000000)
	k.UserID(message.Message.ReplyMessage.FromID)

	_, err := session.MessagesRemoveChatUser(k.Params)
	if err != nil {
		apiutil.Send(session, err.Error(), message.Message.PeerID)
	}
}
