package quote

import (
	"errors"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/beshenkaD/maschinenkonzept/core"
	"github.com/cavaliercoder/grab"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func generateLines(s string) []string {
	const perLine = 6

	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	words := strings.Split(s, " ")

	lines := make([]string, len(words)/(perLine+1)+1)

	current := 0
	line := 0
	for _, word := range words {
		if current == perLine {
			current = 0
			line++
			lines[line] += word + " "
		} else {
			lines[line] += word + " "
			current++
		}
	}

	lines[0] = "«" + lines[0]
	lines[len(lines)-1] = strings.TrimSpace(lines[len(lines)-1]) + "»"

	return lines
}

func calculateHeight(fontSize, spacing, linesCount int) int {
	const min = 400

	h := ((fontSize + 5) * linesCount) + 5

	if h < min {
		return min
	}

	return h
}

func calculateStringWidth(s string, fontSize int) int {
	return len(s) * fontSize
}

func generateQuote(photo image.Image, name, quote string, bg, fg color.Color) string {
	lines := generateLines(quote)

	const fontSize = 20
	const spacing = 20

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: fontSize})

	const W = 800
	H := calculateHeight(fontSize, spacing, len(lines))

	dc := gg.NewContext(W, H)

	dc.SetFontFace(face)
	dc.SetColor(bg)
	dc.Clear()
	dc.SetColor(fg)

	for i, line := range lines {
		y := H/2 - spacing*len(lines)/2 + i*spacing
		dc.DrawStringAnchored(line, W/3+W/3, float64(y), 0.5, 0.5)
	}

	dc.DrawStringAnchored(name, W/7, float64(H-30), 0.5, 0.5)

	dc.DrawEllipse(W/7, float64(H/2), 100, 100)
	dc.Clip()
	dc.DrawImageAnchored(photo, W/7, H/2, 0.5, 0.5)

	path := filepath.Join(os.TempDir(), strconv.Itoa(int(time.Now().UnixNano()))+".png")
	dc.SavePNG(path)

	return path
}

func quote(i *core.CommandInput) (string, error) {
	var name string
	var photoURL string

	if i.Message.ReplyMessage == nil {
		return "", errors.New("вы не ответили ни на какое сообщение")
	}

	if i.Message.ReplyMessage.FromID < 0 {
		bot, err := core.Vk.GroupsGetByID(api.Params{
			"group_ids": -i.Message.ReplyMessage.FromID,
			"lang":      0,
			"fields":    "photo_200",
		})

		if err != nil {
			return "", err
		}

		name = bot[0].Name + " ©"
		photoURL = bot[0].Photo200
	} else {
		user, err := core.Vk.UsersGet(api.Params{
			"user_ids": i.Message.ReplyMessage.FromID,
			"lang":     0,
			"fields":   "photo_200",
		})

		if err != nil {
			return "", err
		}

		name = user[0].FirstName + " " + user[0].LastName + " ©"

		if i.User.ID == i.Message.ReplyMessage.FromID {
			name += " (Self)"
		}

		photoURL = user[0].Photo200
	}

	resp, err := grab.Get(os.TempDir(), photoURL)

	if err != nil {
		return "", err
	}

	photo, err := gg.LoadImage(resp.Filename)
	if err != nil {
		return "", err
	}

	var (
		fg color.Color = color.White
		bg color.Color = color.Black
	)

	if len(i.Args) > 0 {
		switch i.Args[0] {
		case "dark":
		case "light":
			bg = color.White
			fg = color.Black
		default:
			return "", errors.New("f")
		}
	}

	quotePath := generateQuote(photo, name, i.Message.ReplyMessage.Text, bg, fg)

	err = core.SendMessage(i.Chat, "Вот ваша цитата", "", quotePath, nil)

	os.Remove(quotePath)

	return "", err
}

func init() {
	rand.Seed(time.Now().UnixNano())
	core.RegisterCommand("quote", "", nil, quote)
}
