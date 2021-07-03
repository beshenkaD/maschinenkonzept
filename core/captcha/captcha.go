package captcha

import (
	"container/heap"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/beshenkaD/maschinenkonzept/apiutil"
	"github.com/beshenkaD/maschinenkonzept/core"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type userTimeout struct {
	ID      int
	ChatID  int
	Answer  int
	Start   time.Time
	Correct bool
}

type userTimeoutHeap []userTimeout

func (h userTimeoutHeap) Len() int           { return len(h) }
func (h userTimeoutHeap) Peek() *userTimeout { return &h[0] }
func (h userTimeoutHeap) Less(i, j int) bool { return h[i].Start.Before(h[j].Start) }
func (h userTimeoutHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *userTimeoutHeap) Push(x interface{}) {
	*h = append(*h, x.(userTimeout))
}

func (h *userTimeoutHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

// Предлагает пройти капчу только что вступившему в беседу человеку
type CaptchaModule struct {
	timeouts    *userTimeoutHeap
	timeoutLock sync.Mutex
}

func New() *CaptchaModule {
	w := &CaptchaModule{
		timeouts: &userTimeoutHeap{},
	}

	heap.Init(w.timeouts)

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

	println(first, second, answer, first+second)

	s := fmt.Sprintf("Пожалуйста решите пример %d + %d", first, second)
	apiutil.Send(bot.Session, s, msg.Message.PeerID)

	timeout := userTimeout{
		ID:     ID,
		ChatID: msg.Message.PeerID - 2000000000,
		Start:  time.Now(),
		Answer: answer,
	}

	w.timeouts.Push(timeout)
}

func (w *CaptchaModule) OnMessage(bot *core.Bot, msg events.MessageNewObject) {
	if w.timeouts.Len() == 0 {
		return
	}

	last := w.timeouts.Peek()

	w.timeoutLock.Lock()
	if msg.Message.FromID == last.ID && msg.Message.PeerID == last.ChatID+2000000000 {
		r, err := strconv.Atoi(msg.Message.Text)
		if err != nil {
			return
		}

		if r == last.Answer {
			w.timeouts.Pop()
		}
	}
	w.timeoutLock.Unlock()
}

func (w *CaptchaModule) OnTick(bot *core.Bot) {
	if w.timeouts.Len() == 0 {
		return
	}

	last := w.timeouts.Peek()

	w.timeoutLock.Lock()

	d := time.Since(last.Start)

	if d.Seconds() >= 30.0 && !last.Correct && w.timeouts.Len() != 0 {
		b := params.NewMessagesRemoveChatUserBuilder()
		b.ChatID(last.ChatID)
		b.MemberID(last.ID)

		_, err := bot.Session.MessagesRemoveChatUser(b.Params)
		if err != nil {
			println(err.Error())
		}

		w.timeouts.Pop()
	}

	w.timeoutLock.Unlock()
}

func generateCaptcha() (int, int, int) {
	rand.Seed(time.Now().UnixNano())

	answer := rand.Intn(30)
	first := rand.Intn(answer)
	second := answer - first

	return first, second, answer
}
