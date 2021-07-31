package main

import (
	"context"
	"log"
	"os"

	"github.com/beshenkaD/maschinenkonzept/core"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

var (
	vk  *api.VK
	bot *core.Bot
)

func test(in *core.Input) (string, error) {
	return "Ну да я работаю ок", nil
}

func testHook(in *core.Input) (string, error) {
	return "Ты ахуел? обратно смени", nil
}

func testDisable(in *core.Input) (string, error) {
	bot.DisableCommand("test", in.Chat)

	return "я выключил команду test", nil
}

func testEnable(in *core.Input) (string, error) {
	bot.EnableCommand("test", in.Chat)

	return "я включил команду test", nil
}

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

	bot = core.New(responseHandler, "vk")

	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		target := obj.Message.PeerID

		message := &core.Message{
			Text:       obj.Message.Text,
			ActionType: obj.Message.Action.Type,
			IsPrivate:  obj.Message.PeerID < 2000000000,
		}

		from := obj.Message.FromID

		bot.MessageReceived(target, message, &core.User{
			ID:    from,
			Name:  "TODO",
			IsBot: from < 0,
		})

	})

	log.Println("Start Long Poll (VK)")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	core.RegisterHook("a", "chat_title_update", "fd", testHook)
	core.RegisterCommand("test", "test", "test", test)
	core.RegisterCommand("test", "выкл", "", testDisable)
	core.RegisterCommand("test", "вкл", "", testEnable)

	Run(os.Getenv("VK_TOKEN"), true)
}
