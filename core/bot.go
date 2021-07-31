package core

type ResponseHandler func(chat int, message string)

// TODO: cron for periodic commands
type Bot struct {
	handler ResponseHandler

	Protocol string // vk or telegram

	CmdPrefix        map[int]string
	DisabledCommands map[int][]string
	DisabledHooks    map[int][]string

	done chan struct{}
}

func New(h ResponseHandler, protocol string) *Bot {
	b := &Bot{
		handler:          h,
		Protocol:         protocol,
		done:             make(chan struct{}),
		CmdPrefix:        make(map[int]string),
		DisabledCommands: make(map[int][]string),
		DisabledHooks:    make(map[int][]string),
	}

	return b
}

func (b *Bot) ChangePrefix(prefix string, chat int) {
	b.CmdPrefix[chat] = prefix[1:]
}

func (b *Bot) startPeriodic() {}

func (b *Bot) MessageReceived(chat int, message *Message, sender *User) {
	prefix, ok := b.CmdPrefix[chat]

	if !ok {
		prefix = "/"
	}

	input, err := parse(message, chat, sender, prefix)

	if err != nil {
		b.SendMessage(chat, err.Error())
		return
	}

	if input == nil {
		return
	}

	if b.IsCommandDisabled(input.Command, chat) {
		return
	}

	go b.handleCmd(input)
	go b.handleHook(input)
}

func (b *Bot) SendMessage(chat int, message string) {
	if message == "" {
		return
	}

	b.handler(chat, message)
}

func (b *Bot) Stop() {
	close(b.done)
}
