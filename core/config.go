package core

import (
	"reflect"
)

type moduleID string
type commandID string

type Config struct {
	SetupDone bool
	Basic     struct {
		IgnoreInvalidCommands bool
		Aliases               map[string]string
		CommandPrefix         string
	}

	Modules struct {
		Disabled        map[moduleID]bool
		CommandDisabled map[commandID]bool
	}
}

func DefaultConfig() *Config {
	config := &Config{
		SetupDone: false,
	}

	config.Basic.CommandPrefix = "/"
	config.Basic.IgnoreInvalidCommands = false

	return config
}

func setConfigValue(f reflect.Value, value string, chat *Chat) error {
	switch f.Interface().(type) {
	case string:
		if value == `""` {
			f.SetString("")
		} else {
			f.SetString(value)
		}
	}

	return nil
}
