package core

func (b *Bot) OnMessage(chat *Chat, msg vkMessage) {
	b.Processed++
	pm := chat.ID < 2000000000

	go runCommand(msg, chat, pm)
}

func (b *Bot) OnPhotoUpdate(chat *Chat, msg vkMessage) {
	for _, h := range chat.hooks.OnPhotoUpdate {
		if chat.ShouldRunHooks(h) {
			go h.OnPhotoUpdate(chat, msg)
		}
	}
}

func (b *Bot) OnPhotoRemove(chat *Chat, msg vkMessage) {
	for _, h := range chat.hooks.OnPhotoRemove {
		if chat.ShouldRunHooks(h) {
			go h.OnPhotoRemove(chat, msg)
		}
	}
}

func (b *Bot) OnChatCreate(chat *Chat, msg vkMessage) {
	for _, h := range chat.hooks.OnChatCreate {
		if chat.ShouldRunHooks(h) {
			go h.OnChatCreate(chat, msg)
		}
	}
}

func (b *Bot) OnTitleUpdate(chat *Chat, msg vkMessage) {
	for _, h := range chat.hooks.OnTitleUpdate {
		if chat.ShouldRunHooks(h) {
			go h.OnTitleUpdate(chat, msg)
		}
	}
}

func (b *Bot) OnInviteUser(chat *Chat, msg vkMessage) {
	for _, h := range chat.hooks.OnInviteUser {
		if chat.ShouldRunHooks(h) {
			go h.OnInviteUser(chat, msg)
		}
	}
}

func (b *Bot) OnKickUser(chat *Chat, msg vkMessage) {
	for _, h := range chat.hooks.OnKickUser {
		if chat.ShouldRunHooks(h) {
			go h.OnKickUser(chat, msg)
		}
	}
}

func (b *Bot) OnPinMessage(chat *Chat, msg vkMessage) {
	for _, h := range chat.hooks.OnPinMessage {
		if chat.ShouldRunHooks(h) {
			go h.OnPinMessage(chat, msg)
		}
	}
}

func (b *Bot) OnUnpinMessage(chat *Chat, msg vkMessage) {
	for _, h := range chat.hooks.OnUnpinMessage {
		if chat.ShouldRunHooks(h) {
			go h.OnUnpinMessage(chat, msg)
		}
	}
}

func (b *Bot) OnInviteByLink(chat *Chat, msg vkMessage) {
	for _, h := range chat.hooks.OnInviteByLink {
		if chat.ShouldRunHooks(h) {
			go h.OnInviteByLink(chat, msg)
		}
	}
}
