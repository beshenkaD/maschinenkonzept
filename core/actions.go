package core

type actionType string

const (
	ActionNewMessage   = "chat_message_new"
	ActionPhotoUpdate  = "chat_photo_update"
	ActionPhotoRemove  = "chat_photo_remove"
	ActionChatCreate   = "chat_create"
	ActionTitleUpdate  = "chat_title_update"
	ActionInviteUser   = "chat_invite_user"
	ActionKickUser     = "chat_kick_user"
	ActionPinMessage   = "chat_pin_message"
	ActionUnpinMessage = "chat_unpin_message"
	ActionInviteByLink = "chat_invite_user_by_link"
)

func parseAction(action string) actionType {
	if action == "" {
		return ActionNewMessage
	}

	return actionType(action)
}
