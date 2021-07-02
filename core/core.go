package core

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
)

func (b *Bot) Run() {
	group, err := b.Session.GroupsGetByID(nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Запущен как: ", group[0].Name)

	lp, err := longpoll.NewLongPoll(b.Session, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)

        b.RunListeners(obj)

		if cmdFunc, ok := b.Commands[obj.Message.Text]; ok {
			go cmdFunc(b.Session, obj)
		}
	})

	log.Println("Start Long Poll")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}
