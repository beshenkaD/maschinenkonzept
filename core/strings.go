package core

type Strings struct {
	InvalidCommand  string
	DisabledCommand string
}

func Lang(chat int) *Strings {
	lang, ok := language[chat]

	if !ok {
		return StringMap[DefaultLanguage]
	}

	return StringMap[lang]
}

var StringMap = map[string]*Strings{
	"ru": {
		InvalidCommand:  "Неправильная команда",
		DisabledCommand: "Команда отключена в данной беседе",
	},
	"en": {
		InvalidCommand:  "Invalid command",
		DisabledCommand: "Command disabled",
	},
	"ua": {
		InvalidCommand:  "Неправильна команда",
		DisabledCommand: "хрю хрю хрю",
	},
}
