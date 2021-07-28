package core

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type moduleID string
type commandID string

type Config struct {
	SetupDone bool `json:"setupdone"`
	Basic     struct {
		IgnoreInvalidCommands bool              `json:"ignoreinvalidcommands"`
		Aliases               map[string]string `json:"aliases"`
		CommandPrefix         string            `json:"commandprefix"`
	} `json:"basic"`

	Modules struct {
		Disabled         map[moduleID]bool  `json:"disabledmodules"`
		DisabledCommands map[commandID]bool `json:"disabledcommands"`
	} `json:"modules"`
}

func NewConfig() *Config {
	config := &Config{
		SetupDone: false,
	}

	config.Basic.CommandPrefix = "/"
	config.Basic.IgnoreInvalidCommands = false

	return config
}

func (c *Config) Set(chat *Chat, args []string, message string) (string, bool) {
	name := args[0]
	names := strings.SplitN(strings.ToLower(name), ".", 3)

	t := reflect.ValueOf(c).Elem()

	for i := 0; i < t.NumField(); i++ {
		if strings.ToLower(t.Type().Field(i).Name) == names[0] {
			if len(names) < 2 {
				return "err", false
			}
			switch t.Field(i).Kind() {
			case reflect.Struct:
				for j := 0; j < t.Field(i).NumField(); j++ {
					if strings.ToLower(t.Field(i).Type().Field(j).Name) == names[1] {
						f := t.Field(i).Field(j)
						switch f.Interface().(type) {
						case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, float32, float64, uint64:
							value := args[1]

							if err := setConfigValue(f, value, chat); err != nil {
								return "Error: " + err.Error(), false
							}

						case bool:
							switch strings.ToLower(args[1]) {
							case "true", "да":
								f.SetBool(true)
							case "false", "нет":
								f.SetBool(false)
							default:
								return name + " must be set to either 'true' or 'false'", false
							}

						case map[moduleID]bool, map[commandID]bool:
							return setConfigList(f, args[1:], chat)

						case map[string]string:
							value := args[2]

							return setConfigMap(f, strings.ToLower(args[1]), value, chat)

						default:
							return "unknown type", false
						}
						return fmt.Sprint(f.Interface()), true
					}
				}
			default:
				return "not a category", false
			}
		}
	}

	return "не могу найти параметр: " + name, false
}

func setConfigValue(f reflect.Value, value string, chat *Chat) error {
	switch f.Interface().(type) {
	case string:
		if value == `""` {
			f.SetString("")
		} else {
			f.SetString(value)
		}
	case moduleID:
		value := strings.ToLower(value)
		for _, m := range chat.Modules {
			if value == strings.ToLower(m.Name()) {
				f.SetString(value)
				return nil
			}
		}
		return fmt.Errorf("%s not a module", value)
	case commandID:
		value := strings.ToLower(value)
		if _, ok := chat.commands[commandID(value)]; !ok {
			return fmt.Errorf("%s not a command", value)
		}
		f.SetString(value)
	case int, int8, int16, int32, int64:
		k, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		f.SetInt(k)
	case uint, uint8, uint16, uint32, uint64:
		k, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		f.SetUint(k)
	case float32, float64:
		k, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		f.SetFloat(k)
	}

	return nil
}

func setConfigList(f reflect.Value, values []string, chat *Chat) (string, bool) {
	switch f.Kind() {
	case reflect.Slice:
		f.Set(reflect.MakeSlice(f.Type(), 0, len(values)))
		if len(values) > 0 && len(values[0]) > 0 {
			for _, value := range values {
				v := reflect.New(f.Type().Elem()).Elem()
				if err := setConfigValue(v, value, chat); err != nil {
					return "Value error: " + err.Error(), false
				}
				f.Set(reflect.Append(f, v))
			}
		}
		return fmt.Sprint(f.Interface()), true
	case reflect.Map:
		if f.Type().Elem() != reflect.TypeOf(true) {
			return "Map sent into list function!", false
		}
		f.Set(reflect.MakeMap(f.Type()))
		stripped := []string{}
		if len(values) > 0 && len(values[0]) > 0 {
			for _, value := range values {
				v := reflect.New(f.Type().Key()).Elem()
				if err := setConfigValue(v, value, chat); err != nil {
					return "Value error: " + err.Error(), false
				}
				f.SetMapIndex(v, reflect.ValueOf(true))
				stripped = append(stripped, fmt.Sprint(v.Interface()))
			}
		}
		return "[" + strings.Join(stripped, ", ") + "]", true
	}
	return "Unknown list type!", false
}

func deleteFromMapReflect(f reflect.Value, k reflect.Value) string {
	if (f.MapIndex(k) == reflect.Value{}) {
		return fmt.Sprint(k.Interface()) + " does not exist."
	}
	f.SetMapIndex(k, reflect.Value{})
	return "Deleted " + fmt.Sprint(k.Interface())
}

func setConfigMap(f reflect.Value, key, value string, chat *Chat) (string, bool) {
	k := reflect.New(f.Type().Key()).Elem()

	if err := setConfigValue(k, key, chat); err != nil {
		return "Key error: " + err.Error(), false
	}

	if f.IsNil() {
		f.Set(reflect.MakeMap(f.Type()))
	}

	if len(value) == 0 {
		return deleteFromMapReflect(f, k), false
	}

	v := reflect.New(f.Type().Elem()).Elem()
	if err := setConfigValue(v, value, chat); err != nil {
		return "Value error: " + err.Error(), false
	}

	f.SetMapIndex(k, v)
	return fmt.Sprintf("%v: %v", k.Interface(), v.Interface()), true
}
