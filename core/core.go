package core

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
)

// У каждого бота будет свой лонгпол и список плагинов
func Run(bot Bot) {
	bot.Register("ping", func(session *api.VK, message events.MessageNewObject) {
		b := params.NewMessagesSendBuilder()
		b.Message("pong")
		b.RandomID(0)
		b.PeerID(message.Message.PeerID)

		_, err := session.MessagesSend(b.Params)

		if err != nil {
			log.Fatal(err)
		}
	})

	bot.Register("убирай пинг", func(session *api.VK, message events.MessageNewObject) {
		b := params.NewMessagesSendBuilder()
        b.Message("Ок убираю пинг :))")
		b.RandomID(0)
		b.PeerID(message.Message.PeerID)

		_, err := session.MessagesSend(b.Params)

		if err != nil {
			log.Fatal(err)
		}

        bot.Unregister("ping")
	})


	vk := api.NewVK(bot.Token)

	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Запущен как: ", group[0].Name)

	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	// New message event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)

		if cmdFunc, ok := bot.Commands[obj.Message.Text]; ok {
			cmdFunc(vk, obj)
		}
	})

	// Run Bots Long Poll
	log.Println("Start Long Poll")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}
