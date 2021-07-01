package core

/// Дефолтные команды бота

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"runtime"
)

func ping(session *api.VK, message events.MessageNewObject) {
	b := params.NewMessagesSendBuilder()
	b.Message("pong")
	b.RandomID(0)
	b.PeerID(message.Message.PeerID)

	_, err := session.MessagesSend(b.Params)

	if err != nil {

	}
}

func stat(session *api.VK, message events.MessageNewObject) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	b := params.NewMessagesSendBuilder()
	b.Message(fmt.Sprintf(`
Total alloc: %v MiB
System: %v Mib`, bToMb(m.TotalAlloc), bToMb(m.Sys)))

	b.RandomID(0)
	b.PeerID(message.Message.PeerID)

	_, err := session.MessagesSend(b.Params)

	if err != nil {
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
