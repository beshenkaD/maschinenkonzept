package vkutil

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
)

func RemoveUser(session *api.VK, chatID, userID int) (int, error) {
    b := params.NewMessagesRemoveChatUserBuilder()
    b.ChatID(chatID)
    b.UserID(userID)

    r, err := session.MessagesRemoveChatUser(b.Params)

    return r, err
}
