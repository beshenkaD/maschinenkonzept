package quote

import (
	"errors"
	"image"
	"image/color"
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
)

const (
	fontSize = 20
	// Image must always be 700 pixels wide, but height may vary
	width     = 700
	minHeight = 400
)

var (
	face = getFontFace(fontSize, true)
)

func getFontFace(size int, regular bool) font.Face {
	var font *truetype.Font

	// TODO: find it automatically and set it in init func
	r, _ := os.ReadFile("/usr/share/fonts/droid/DroidSansFallbackFull.ttf")
	b, _ := os.ReadFile("/usr/share/fonts/droid/DroidSans-Bold.ttf")
	if regular {
		font, _ = truetype.Parse(r)
	} else {
		font, _ = truetype.Parse(b)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: float64(size)})

	return face
}

func getPhotoPoint(height int) (x, y int) {
	return width / 6, height / 2
}

func getNamePoint(height int) (x, y int) {
	return 15, height - 15
}

func getStringWidth(s string) int {
	w := 0
	for _, r := range s {
		_, a, _ := face.GlyphBounds(r)
		w += a.Round()
	}

	return w
}

func getLinesHeight(lines []string) int {
	h := (len(lines) * (face.Metrics().Height.Ceil() + 2))
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
			if getStringWidth(newLine+" "+word) > (width - w - 10) {
				lines = append(lines, newLine)
				newLine = word
			} else {
				newLine = newLine + " " + word
			}
		}

		newLine = strings.TrimSpace(newLine)
		if len(newLine) != 0 {
			lines = append(lines, newLine)
		}
		newLine = ""
	}

	lines[0] = "«" + lines[0]
	lines[len(lines)-1] = lines[len(lines)-1] + "»"

	return lines
}

func getRainbowName(firstName, lastName string, self bool) string {
	sub := map[string]string{
		"яна":  "лесбияна",
		"гей":  "gay",
		"ге":   "gay",
		"де":   "gay",
		"пере": "пидор",
		"слав": "slave",
		"кам":  "cum",
	}

	for k, v := range sub {
		firstName = strings.ReplaceAll(strings.ToLower(firstName), k, v)
		lastName = strings.ReplaceAll(strings.ToLower(lastName), k, v)
	}

	return getName(firstName, lastName, self)
}

func getName(firstName, lastName string, self bool) string {
	s := ""

	if self {
		s = "(Self Signed)"
	}

	return strings.Title(firstName) + " " + strings.Title(lastName) + " (c) " + s
}

type quoteMode int

const (
	lightMode = iota
	darkMode
	rainbowMode
)

func generateQuote(photo image.Image, firstName, lastName, quote string, self bool, mode quoteMode) (path string) {
	var (
		fg     color.Color
		bg     color.Color
		name   string
		drawer func()
	)

	const (
		textPointX = width / 3
	)

	lines := getLines(quote, textPointX)
	height := getLinesHeight(lines)

	dc := gg.NewContext(width, height)

	// Drawer for rainbow mode
	rainbow := func() {
		colors := []color.Color{
			// Red
			color.RGBA{255, 0, 0, 255},
			// Orange
			color.RGBA{255, 128, 0, 255},
			// Yellow
			color.RGBA{255, 186, 0, 255},
			// Green
			color.RGBA{0, 255, 0, 255},
			// Light blue
			color.RGBA{0, 255, 255, 255},
			// Blue
			color.RGBA{0, 0, 255, 255},
			// Violet
			color.RGBA{255, 0, 255, 255},
		}

		bold := getFontFace(35, false)
		h := height + bold.Metrics().Height.Ceil() + 10

		dc.SetFontFace(bold)
		dc.DrawStringAnchored("ЦИТАТЫ ВЕЛИКИХ ПИДОРАСОВ", width/2, 20, 0.5, 0.5)
		dc.SetFontFace(face)

		c := 0
		for i, line := range lines {
			y := h/2 - fontSize*len(lines)/2 + i*fontSize

			x := textPointX
			for _, word := range strings.Split(line, " ") {
				for _, r := range word {
					dc.SetColor(colors[c])

					if c == 6 {
						c = 0
					}

					dc.DrawString(string(r), float64(x), float64(y))
					x += getStringWidth(string(r))
					c++
				}
				x += getStringWidth(" ")
			}
		}
		dc.SetColor(fg)
	}

	classic := func() {
		for i, line := range lines {
			y := height/2 - fontSize*len(lines)/2 + i*fontSize
			dc.DrawString(line, textPointX, float64(y))
		}
	}

	switch mode {
	case lightMode:
		fg = color.Black
		bg = color.White
		name = getName(firstName, lastName, self)
		drawer = classic
	case darkMode:
		fg = color.White
		bg = color.Black
		name = getName(firstName, lastName, self)
		drawer = classic
	case rainbowMode:
		fg = color.Black
		bg = color.White
		name = getRainbowName(firstName, lastName, self)
		drawer = rainbow
	default:
		fg = color.White
		bg = color.Black
		name = getName(firstName, lastName, self)
		drawer = classic
	}

	dc.SetFontFace(face)
	dc.SetColor(bg)
	dc.Clear()
	dc.SetColor(fg)

	// Draw quote text
	drawer()

	// Draw name
	nx, ny := getNamePoint(height)
	dc.DrawString(name, float64(nx), float64(ny))

	// Draw time
	t := time.Now().UTC().Format("02.01.2006 15:04")
	dc.DrawString(t, float64(nx+510), float64(ny))

	// Draw photo and make it round
	px, py := getPhotoPoint(height)
	dc.DrawEllipse(float64(px), float64(py), 100, 100)
	dc.Clip()
	dc.DrawImageAnchored(photo, px, py, 0.5, 0.5)

	out := filepath.Join(os.TempDir(), time.Now().String()+`.png`)
	dc.SavePNG(out)

	return out
}

func getUserInfo(ID int) (firstName, lastName, photoURL string, err error) {
	if ID < 0 {
		bot, err := core.Vk.GroupsGetByID(api.Params{
			"group_ids": -ID,
			"lang":      0,
			"fields":    "photo_200",
		})

		if err != nil {
			return "", "", "", err
		}

		return bot[0].Name, "", bot[0].Photo200, nil
	} else {
		user, err := core.Vk.UsersGet(api.Params{
			"user_ids": ID,
			"lang":     0,
			"fields":   "photo_200",
		})

		if err != nil {
			return "", "", "", err
		}

		return user[0].FirstName, user[0].LastName, user[0].Photo200, nil
	}
}

func quote(i *core.CommandInput) (string, error) {
	var (
		firstName string
		lastName  string
		photoURL  string
		text      string
		self      bool
		err       error
	)

	if i.Message.ReplyMessage != nil {
		firstName, lastName, photoURL, err = getUserInfo(i.Message.ReplyMessage.FromID)
		if err != nil {
			return "", err
		}
		self = i.User.ID == i.Message.ReplyMessage.FromID

		text = i.Message.ReplyMessage.Text
	} else if len(i.Message.FwdMessages) != 0 {
		ID := i.Message.FwdMessages[0].FromID

		text = ""
		for _, msg := range i.Message.FwdMessages {
			if ID == msg.FromID {
				text += msg.Text + "\n"
			}
		}

		firstName, lastName, photoURL, err = getUserInfo(ID)
		if err != nil {
			return "", err
		}
		self = ID == i.User.ID
	} else {
		return "ПОШЕЛ НАХУЙ", nil
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
		mode quoteMode
	)

	if len(i.Args) > 0 {
		switch i.Args[0] {
		case "dark":
			mode = darkMode
		case "light":
			mode = lightMode
		case "rainbow":
			mode = rainbowMode
		default:
			return "", errors.New("a")
		}
	}

	quotePath := generateQuote(photo, firstName, lastName, text, self, mode)

	err = core.SendMessage(i.Chat, "", "", quotePath, nil)

	os.Remove(quotePath)

	return "", err
}

func init() {
	core.RegisterCommand("quote", "", nil, quote)
}
