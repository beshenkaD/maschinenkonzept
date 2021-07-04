package captcha

import (
	// "container/heap"
	"fmt"
	// "github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/beshenkaD/maschinenkonzept/apiutil"
	"github.com/beshenkaD/maschinenkonzept/core"
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

// Предлагает пройти капчу только что вступившему в беседу человеку
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

func (w *CaptchaModule) OnInviteUser(bot *core.Bot, msg events.MessageNewObject) {
	ID := msg.Message.Action.MemberID

	if ID < 0 {
		return
	}

	first, second, answer := generateCaptcha()

    name := apiutil.GetName(bot.Session, ID)
	s := fmt.Sprintf("[id%d|%s], Пожалуйста решите пример %d + %d", ID, name, first, second)
	apiutil.Send(bot.Session, s, msg.Message.PeerID)

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

			apiutil.Send(bot.Session, fmt.Sprintf("%d ты бот ёбаный", ID), timeout.chat+2000000000)
		}
	}
}

func generateCaptcha() (int, int, int) {
	rand.Seed(time.Now().UnixNano())

    min := 5
    max := 30

	answer := rand.Intn(max - min + 1) + min
	first := rand.Intn(answer)
	second := answer - first

	return first, second, answer
}
