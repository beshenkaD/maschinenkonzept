package core

import (
	"math/rand"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
)

func SendMessage(chat *Chat, message string, attachment string, attachmentPath string, replyTo *Message) error {
	b := params.NewMessagesSendBuilder()
	b.Lang(0)
	b.PeerID(chat.ID)
	b.RandomID(int(rand.Int31()))
	b.Message(message)

	if attachment != "" {
		b.Attachment(attachment)
	}

	if attachmentPath != "" {
		// TODO
	}

	if replyTo != nil {
		b.ReplyTo(replyTo.ConversationMessageID)
	}

	_, err := Vk.MessagesSend(b.Params)

	return err
}

func RemoveUser(chat *Chat, userID int) error {
	b := params.NewMessagesRemoveChatUserBuilder()
	b.ChatID(chat.ID - 2000000000)

	if userID > 0 {
		b.UserID(userID)
	} else {
		b.MemberID(userID)
	}

	_, err := Vk.MessagesRemoveChatUser(b.Params)

	return err
}

func GetInviteLink(chat *Chat, reset bool) (string, error) {
	b := params.NewMessagesGetInviteLinkBuilder()
	b.PeerID(chat.ID)
	b.Reset(reset)

	a, err := Vk.MessagesGetInviteLink(b.Params)

	return a.Link, err
}

func DeleteMessages(chat *Chat, messageIds []int) error {
	_, err := Vk.MessagesDelete(api.Params{
		"delete_for_all":           true,
		"peer_id":                  chat.ID,
		"conversation_message_ids": messageIds,
	})

	return err
}

func Unpin(chat *Chat) error {
	return nil
}

func Pin(chat *Chat, message *Message) error {
	return nil
}

func RenameChat(chat *Chat, title string) error {
	return nil
}

func GetChat(chats ...*Chat) (int, string) {
	return 0, ""
}
