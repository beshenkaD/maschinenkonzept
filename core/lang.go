package core

type Languages map[string]*Strings

var (
	langs = make(Languages)
)

type Strings struct {
	CommandDisabled       string
	InvalidCommand        string
	Error                 string
	HelpDescription       string
	HelpUsage             string
	HelpAvailableCommands string
	HelpAboutCommand      string
}

func AddLanguage(code string, strings *Strings) {
	langs[code] = strings
}

func GetAvailableLanguages() []string {
	var l []string

	for lang := range langs {
		l = append(l, lang)
	}

	return l
}

func getStrings(code string) *Strings {
	s, ok := langs[code]

	if !ok {
		return langs["en"]
	}

	return s
}

func init() {
	AddLanguage("en", &Strings{
		CommandDisabled:       "command disabled in this chat",
		InvalidCommand:        "invalid command",
		HelpDescription:       "Description: %s",
		HelpUsage:             "Usage: %s%s %s",
		HelpAvailableCommands: "Available commands: %v",
		HelpAboutCommand:      "Type: '%shelp <command>' to see details about a specific command.",
		Error:                 "Error: ",
	})
}
