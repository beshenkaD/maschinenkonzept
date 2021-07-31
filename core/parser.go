package core

import (
	"errors"
	"strings"

	shellwords "github.com/mattn/go-shellwords"
)

func parse(m *Message, chat int, user *User, prefix string) (*Input, error) {
	s := strings.TrimSpace(m.Text)

	if !strings.HasPrefix(s, prefix) {
		// TODO: Refactor this
		if m.ActionType != "" {
			return &Input{
				Chat:        chat,
				Message:     s,
				User:        user,
				MessageData: m,
			}, nil
		}
		return nil, nil
	}

	i := &Input{
		Chat:    chat,
		Message: strings.TrimSpace(strings.TrimPrefix(s, prefix)),
		User:    user,
	}

	if i.Message == "" {
		return nil, nil
	}

	pieces := strings.Fields(i.Message)
	i.Command = strings.ToLower(pieces[0])

	if len(pieces) > 1 {
		args, err := shellwords.Parse(pieces[1])

		if err != nil {
			return nil, errors.New("error parsing arguments " + err.Error())
		}

		i.Args = args
	}

	m.Text = i.Message
	i.MessageData = m

	return i, nil
}
