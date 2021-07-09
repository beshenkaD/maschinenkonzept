package core

type Config struct {
	SetupDone bool
	Basic     struct {
		IgnoreInvalidCommands bool
		Aliases               map[string]string
		CommandPrefix         string
	}

	Modules struct {
		Disabled        map[string]bool
		CommandDisabled map[string]bool
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
