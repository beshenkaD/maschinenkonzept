package core

import (
	"strings"
	"unicode"
)

func parse(m *Message, chat int, user *User, prefix byte) (*CommandInput, error) {
	s := strings.TrimSpace(m.Text)

	if !strings.HasPrefix(s, string(prefix)) {
		return nil, nil
	}

	i := &CommandInput{
		Chat:    chat,
		Message: strings.TrimSpace(strings.TrimPrefix(s, string(prefix))),
		User:    user,
	}

	if i.Message == "" {
		return nil, nil
	}

	firstOccurrence := true
	firstUnicodeSpace := func(c rune) bool {
		isFirstSpace := unicode.IsSpace(c) && firstOccurrence
		if isFirstSpace {
			firstOccurrence = false
		}
		return isFirstSpace
	}

	pieces := strings.FieldsFunc(i.Message, firstUnicodeSpace)
	i.Command = strings.ToLower(pieces[0])

	if len(pieces) > 1 {
		args := strings.Split(strings.TrimSpace(pieces[1]), " ")

		i.Args = args
	}

	m.Text = i.Message
	i.MessageData = m

	return i, nil
}
