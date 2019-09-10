package scoreboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/iDigitalFlame/scorebot-scoreboard/scoreboard/web"
)

// Log is a struct that stores and represents the Scoreboard Logging config
// Able to be loaded from JSON
type Log struct {
	File  string `json:"file"`
	Level uint8  `json:"level"`
}

// Config is a struct that stores and represents the Scoreboard config
// Able to be loaded from JSON.
type Config struct {
	Log       *Log     `json:"log,omitempty"`
	Key       string   `json:"key"`
	Cert      string   `json:"cert"`
	Tick      uint16   `json:"tick"`
	Assets    string   `json:"assets"`
	Listen    string   `json:"listen"`
	Twitter   *Twitter `json:"twitter,omitempty"`
	Timeout   uint16   `json:"timeout"`
	Scorebot  string   `json:"scorebot"`
	Directory string   `json:"dir"`
}

// Twitter is a struct that stores and represents the Scoreboard Twitter config
// Able to be loaded from JSON.
type Twitter struct {
	Filter      *web.Filter      `json:"filter"`
	Expire      uint16           `json:"expire"`
	Timeout     uint16           `json:"timeout"`
	Credentials *web.Credentials `json:"auth"`
}

// Defaults returns a JSON string representation of the default config.
// Used for creating and understanding the config file structure.
func Defaults() string {
	c := &Config{
		Log: &Log{
			File:  "",
			Level: DefaultLogLevel,
		},
		Key:    "",
		Cert:   "",
		Tick:   DefaultTick,
		Assets: "",
		Listen: DefaultListen,
		Twitter: &Twitter{
			Filter: &web.Filter{
				Language: []string{"en"},
				Keywords: []string{
					"pvj",
					"ctf",
				},
				OnlyUsers:    []string{},
				BlockedUsers: []string{},
				BlockedWords: []string{},
			},
			Expire:  DefaultExpire,
			Timeout: DefaultTimeout,
			Credentials: &web.Credentials{
				AccessKey:      "",
				ConsumerKey:    "",
				AccessSecret:   "",
				ConsumerSecret: "",
			},
		},
		Timeout:   DefaultTimeout,
		Scorebot:  "http://scorebot",
		Directory: "html",
	}
	b, _ := json.MarshalIndent(c, "", "    ")
	return string(b)
}
func (c *Config) verify() error {
	if c.Tick <= 0 {
		return ErrInvalidTick
	}
	if c.Timeout < 0 {
		return web.ErrInvalidTimeout
	}
	if c.Log != nil {
		if c.Log.Level < 0 || c.Log.Level > 5 {
			return ErrInvalidLevel
		}
	} else {
		c.Log = &Log{Level: DefaultLogLevel}
	}
	if c.Twitter != nil {
		v := true
		if c.Twitter.Credentials != nil {
			if len(c.Twitter.Credentials.AccessKey) == 0 || len(c.Twitter.Credentials.AccessSecret) == 0 {
				v = false
			}
			if len(c.Twitter.Credentials.ConsumerKey) == 0 || len(c.Twitter.Credentials.ConsumerSecret) == 0 {
				v = false
			}
		} else {
			v = false
		}
		if c.Twitter.Filter != nil {
			if len(c.Twitter.Filter.Language) == 0 || len(c.Twitter.Filter.Keywords) == 0 {
				v = false
			}
		} else {
			v = false
		}
		if !v {
			c.Twitter = nil
		}
	}
	return nil
}

// Load loads the config from the specified file path 's'
func Load(s string) (*Config, error) {
	f, err := os.Stat(s)
	if err != nil {
		return nil, fmt.Errorf("cannot load file \"%s\": %w", s, err)
	}
	if f.IsDir() {
		return nil, fmt.Errorf("cannot load \"%s\": path is not a file", s)
	}
	b, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, fmt.Errorf("cannot read file \"%s\": %w", s, err)
	}
	var c *Config
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("cannot read file \"%s\" into JSON: %w", s, err)
	}
	return c, nil
}

// SplitParm returns a string array from a comma seperated list.
// This function also trims the string lengths of excess spaces.
func SplitParm(s, d string) []string {
	if len(s) == 0 {
		return []string{}
	}
	f := strings.Split(s, d)
	for i := range f {
		f[i] = strings.TrimSpace(f[i])
	}
	return f
}