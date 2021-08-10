package quote

import (
	"errors"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/beshenkaD/maschinenkonzept/core"
	"github.com/cavaliercoder/grab"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	// "github.com/golang/freetype/truetype"
	// "golang.org/x/image/font/gofont/goregular"
)

const (
	fontSize = 20
	// Image must always be 800 pixels wide, but height may vary
	width     = 700
	minHeight = 400
)

var (
	face = getFontFace()
)

func getFontFace() font.Face {
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: fontSize})

	return face
}

func getPhotoPoint(height int) (x, y int) {
	return width / 6, height - 200
}

func getNamePoint(height int) (x, y int) {
	return 15, height - 40
}

func getStringWidth(s string) int {
	w := 0
	for _, r := range s {
		_, a, _ := face.GlyphBounds(r)
		w += a.Round()
	}

	return w
}

func getStringHeight() int {
	return face.Metrics().Height.Ceil()
}

func getLinesHeight(lines []string) int {
	h := (len(lines) * (getStringHeight() + 2))
	if h < minHeight {
		return minHeight
	}

	return h
}

func getLines(s string, w int) []string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		var newLine string
		for _, word := range strings.Split(line, " ") {
			if getStringWidth(newLine+" "+word) > (width - w - 20) {
				lines = append(lines, newLine)
				newLine = word
			} else {
				newLine = newLine + " " + word
			}
		}

		newLine = strings.TrimSpace(newLine)
		lines = append(lines, newLine)
		newLine = ""
	}

	lines[0] = "«" + lines[0]
	lines[len(lines)-1] = lines[len(lines)-1] + "»"

	return lines
}

func renderRainbowText(ct *gg.Context, textPoint int, s string) {

}

func GenerateQuote(photo image.Image, name, quote string, bg, fg color.Color) (path string) {
	const textPoint = width / 3

	lines := getLines(quote, textPoint)

	height := getLinesHeight(lines)

	dc := gg.NewContext(width, height)
	dc.SetFontFace(face)
	dc.SetColor(bg)
	dc.Clear()
	dc.SetColor(fg)

	for i, line := range lines {
		y := height/2 - fontSize*len(lines)/2 + i*fontSize
		dc.DrawString(line, textPoint, float64(y))
	}

	nx, ny := getNamePoint(height)
	dc.DrawString(name, float64(nx), float64(ny))

	px, py := getPhotoPoint(height)
	dc.DrawEllipse(float64(px), float64(py), 100, 100)
	dc.Clip()
	dc.DrawImageAnchored(photo, px, py, 0.5, 0.5)

	out := filepath.Join(os.TempDir(), time.Now().String()+`.png`)
	dc.SavePNG(out)

	return out
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

	quotePath := GenerateQuote(photo, name, i.Message.ReplyMessage.Text, bg, fg)

	err = core.SendMessage(i.Chat, "", "", quotePath, nil)

	os.Remove(quotePath)

	return "", err
}

func init() {
	rand.Seed(time.Now().UnixNano())
	core.RegisterCommand("quote", "", nil, quote)
}
