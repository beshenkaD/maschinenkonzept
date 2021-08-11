package core

import (
	"math/rand"
	"os"

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
		file, err := os.Open(attachmentPath)
		if err != nil {
			return err
		}

		defer file.Close()

		photo, err := Vk.UploadMessagesPhoto(chat.ID, file)
		if err != nil {
			return err
		}

		b.Attachment(photo)
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

type item struct {
	MemberID  int
	InvitedBy int
	JoinDate  int
	IsAdmin   bool
}

func GetConversationMembers(chat *Chat) ([]item, error) {
	b := params.NewMessagesGetConversationMembersBuilder()
	b.PeerID(chat.ID)

	r, err := Vk.MessagesGetConversationMembers(b.Params)

	if err != nil {
		return nil, err
	}

	var items []item

	for _, i := range r.Items {
		items = append(items, item{
			MemberID:  i.MemberID,
			InvitedBy: i.InvitedBy,
			JoinDate:  i.JoinDate,
			IsAdmin:   bool(i.IsAdmin),
		})
	}

	return items, nil
}

func IsAdmin(chat *Chat, user *User) (bool, error) {
	items, err := GetConversationMembers(chat)
	if err != nil {
		return false, err
	}

	for _, i := range items {
		if i.IsAdmin && (i.MemberID == user.ID) {
			return true, nil
		}
	}

	return false, nil
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

func GetChat(chat *Chat) (string, error) {
	r, err := Vk.MessagesGetChat(api.Params{
		"chat_id": chat.ID - 2000000000,
	})

	if err != nil {
		return "", err
	}

	return r.Title, nil
}
