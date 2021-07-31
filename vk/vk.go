package vk

import (
	"context"
	"log"
	"strconv"

	"github.com/beshenkaD/maschinenkonzept/core"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/SevereCloud/vksdk/v2/object"
)

var (
	vk *api.VK
)

func responseHandler(chat int, message string) {
	b := params.NewMessagesSendBuilder()
	b.PeerID(chat)
	b.RandomID(0)
	b.Message(message)

	_, err := vk.MessagesSend(b.Params)

	if err != nil {
		log.Println(err.Error())
	}
}

func errorHandler(chat int, er error) {
	b := params.NewMessagesSendBuilder()
	b.PeerID(chat)
	b.RandomID(0)
	b.Message("Ошибка: " + er.Error())

	_, err := vk.MessagesSend(b.Params)

	if err != nil {
		log.Println(err.Error())
	}
}

func getFullName(ID int) (string, string) {
	b := params.NewUsersGetBuilder()

	if ID < 0 {
		return "bot", ""
	}

	b.UserIDs([]string{strconv.Itoa(ID)})
	b.Lang(0)

	users, err := vk.UsersGet(b.Params)

	if err != nil {
		log.Println(err.Error())
		return "", ""
	}

	return users[0].FirstName, users[0].LastName
}

func parseMessage(obj *object.MessagesMessage) *core.Message {
	if obj == nil {
		return nil
	}

	actionType := obj.Action.Type

	if actionType == "" {
		actionType = "message_new"
	}

	return &core.Message{
		Text:       obj.Text,
		ActionType: actionType,
		ActionText: obj.Action.Text,
		MemberId:   obj.Action.MemberID,
		IsPrivate:  obj.PeerID < 2000000000,
	}
}

func parseFwd(obj []object.MessagesMessage) []core.Message {
	var msgs []core.Message

	for _, msg := range obj {
		actionType := msg.Action.Type

		if actionType == "" {
			actionType = "message_new"
		}

		msgs = append(msgs, core.Message{
			Text:         msg.Text,
			ActionType:   actionType,
			MemberId:     msg.Action.MemberID,
			ActionText:   msg.Action.Text,
			FwdMessages:  parseFwd(msg.FwdMessages),
			ReplyMessage: parseMessage(msg.ReplyMessage),
			IsPrivate:    msg.PeerID < 2000000000,
		})
	}

	return msgs
}

func Run(token string, debug bool) {
	vk = api.NewVK(token)

	group, err := vk.GroupsGetByID(nil)

	if err != nil {
		log.Fatal(err)
	}

	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	bot := core.New(responseHandler, errorHandler, "vk")

	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		target := obj.Message.PeerID
		message := parseMessage(&obj.Message)
		from := obj.Message.FromID

		firstName, lastName := getFullName(from)
		bot.MessageReceived(target, message, &core.User{
			ID:        from,
			FirstName: firstName,
			LastName:  lastName,
			IsBot:     from < 0,
		})

	})

	log.Println("Start Long Poll (VK)")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}
