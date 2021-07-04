package captcha

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/beshenkaD/maschinenkonzept/core"
	"github.com/beshenkaD/maschinenkonzept/vkutil"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type userTimeout struct {
	chat   int
	answer int
	time   time.Time
}

type userTimeouts map[int]userTimeout

// Требует пройти капчу если юзер вступил по ссылке
type CaptchaModule struct {
	timeouts    userTimeouts
	timeoutLock sync.Mutex
}

func New() *CaptchaModule {
	w := &CaptchaModule{
		timeouts: make(userTimeouts),
	}

	return w
}

func (w *CaptchaModule) Name() string {
	return "Капча"
}

func (w *CaptchaModule) Description() string {
	return "Капча для отсеивания ботов"
}

func (w *CaptchaModule) Commands() []core.Command {
	return []core.Command{}
}

func (w *CaptchaModule) OnInviteByLink(bot *core.Bot, msg events.MessageNewObject) {
	ID := msg.Message.Action.MemberID
	peerID := msg.Message.PeerID

	if ID < 0 {
		return
	}

	first, second, answer := generateCaptcha()

	user, err := vkutil.GetUser(bot.Session, ID)

	if err != nil {
		vkutil.SendMessage(bot.Session, err.Error(), peerID, true)
		return
	}

	s := fmt.Sprintf("[id%d|%s], пожалуйста, решите пример: %d + %d", ID, user.FirstName, first, second)
	vkutil.SendMessage(bot.Session, s, msg.Message.PeerID, false)

	timeout := userTimeout{
		chat:   msg.Message.PeerID - 2000000000,
		answer: answer,
		time:   time.Now(),
	}

	w.timeoutLock.Lock()

	w.timeouts[ID] = timeout

	w.timeoutLock.Unlock()
}

func (w *CaptchaModule) OnMessage(bot *core.Bot, msg events.MessageNewObject) {
	if timeout, ok := w.timeouts[msg.Message.FromID]; ok {
		a, err := strconv.Atoi(msg.Message.Text)

		if err == nil && timeout.answer == a {
			w.timeoutLock.Lock()

			delete(w.timeouts, msg.Message.FromID)

			w.timeoutLock.Unlock()
		}
	}
}

func (w *CaptchaModule) OnTick(bot *core.Bot) {
	for ID, timeout := range w.timeouts {
		if time.Since(timeout.time).Seconds() >= 30.0 {
			w.timeoutLock.Lock()

			delete(w.timeouts, ID)

			w.timeoutLock.Unlock()

			b := params.NewMessagesRemoveChatUserBuilder()
			b.ChatID(timeout.chat)
			b.UserID(ID)

			_, err := bot.Session.MessagesRemoveChatUser(b.Params)

			if err != nil {
				vkutil.SendMessage(bot.Session, err.Error(), timeout.chat+2000000000, true)
			}
		}
	}
}

func generateCaptcha() (int, int, int) {
	rand.Seed(time.Now().UnixNano())

	min := 5
	max := 30

	answer := rand.Intn(max-min+1) + min
	first := rand.Intn(answer)
	second := answer - first

	return first, second, answer
}
