package core

import (
	"strings"
	"unicode"
)

func parse(m *Message, chat *Chat, user *User) *CommandInput {
	s := strings.TrimSpace(m.Text)

	if !strings.HasPrefix(s, chat.Prefix) {
		return nil
	}

	i := &CommandInput{
		Chat:    chat,
		Message: m,
		User:    user,
	}

	s = strings.TrimSpace(strings.TrimPrefix(s, chat.Prefix))

	if s == "" {
		return nil
	}

	firstOccurrence := true
	firstUnicodeSpace := func(c rune) bool {
		isFirstSpace := unicode.IsSpace(c) && firstOccurrence
		if isFirstSpace {
			firstOccurrence = false
		}
		return isFirstSpace
	}

	pieces := strings.FieldsFunc(s, firstUnicodeSpace)
	i.Command = strings.ToLower(pieces[0])

	if len(pieces) > 1 {
		args := strings.Split(strings.TrimSpace(pieces[1]), " ")

		i.Args = args
	}

	return i
}
