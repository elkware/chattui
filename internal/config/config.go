package config

import (
	"os"
	"strings"
)

var AppConfig = NewConfig().Load()

type Config struct {
	ApiKey             string
	CustomInstructions string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	cfg, err := os.ReadFile(homeDir + "/.chattui.conf")
	if err != nil {
		if os.IsNotExist(err) {
			// Return an empty config
			c.ApiKey = ""
			c.CustomInstructions = ""
			return c
		} else {
			panic(err)
		}
	}
	cfgString := strings.Split(string(cfg), "\n")
	for i := range cfgString {
		if strings.HasPrefix(cfgString[i], "API_KEY=") {
			c.ApiKey = cfgString[i][strings.Index(cfgString[i], "API_KEY=")+8:]
			continue
		}
		if strings.HasPrefix(cfgString[i], "CUSTOM_INSTRUCTIONS=") {
			c.CustomInstructions = cfgString[i][strings.Index(cfgString[i], "CUSTOM_INSTRUCTIONS=")+20:]
			continue
		}
	}
	return c
}

func (c *Config) Save() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(homeDir+"/.chattui.conf", []byte("API_KEY="+c.ApiKey+"\nCUSTOM_INSTRUCTIONS="+c.CustomInstructions), 0644)
	if err != nil {
		panic(err)
	}
	return c
}
