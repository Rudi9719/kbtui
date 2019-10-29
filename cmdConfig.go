// +build !rm_basic_commands allcommands setcmd

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml"
)

func init() {
	command := Command{
		Cmd:         []string{"config"},
		Description: "Change various settings",
		Help:        "",
		Exec:        cmdConfig,
	}

	RegisterCommand(command)
}

func cmdConfig(cmd []string) {
	var err error
	switch {
	case len(cmd) == 2:
		if cmd[1] == "load" {
			config, err = readConfig()
			if err != nil {
				printError(err.Error())
				return
			}
			printInfoF("Config file loaded: $TEXT", config.Colors.Feed.File.stylize(config.filepath))
			return
		}
	case len(cmd) > 2:
		if cmd[1] == "load" {
			config, err = readConfig(cmd[3])
			if err != nil {
				printError(err.Error())
				return
			}
			printInfoF("Config file loaded: $TEXT", config.Colors.Feed.File.stylize(config.filepath))
			return
		}
	}
	printError("Must pass a valid command")
}

func readConfig(filepath ...string) (*Config, error) {
	var result = new(Config)
	var configFile string
	var env bool

	// Load default config first, this way any values missing from the provided config file will remain the default value
	d := []byte(defaultConfig)
	toml.Unmarshal(d, result)

	switch len(filepath) {
	case 0:
		configFile, env = os.LookupEnv("KBTUI_CFG")
		if !env {
			configFile = "~/.config/kbtui.toml"
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				configFile = "kbtui.toml"
			}
		}
	default:
		configFile = filepath[0]
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			return result, fmt.Errorf("Unable to load config: %s not found", configFile)
		}
	}

	f, err := ioutil.ReadFile(configFile)
	if err != nil {
		f = []byte(defaultConfig)
	}

	err = toml.Unmarshal(f, result)
	if err != nil {
		return result, err
	}

	result.filepath = configFile
	return result, nil
}
