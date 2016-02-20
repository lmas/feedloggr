package feedloggr2

import (
	"encoding/json"
	"io/ioutil"
)

type FeedConfig struct {
	Title string
	URL   string
}

type Config struct {
	Verbose    bool
	Database   string
	OutputPath string
	Feeds      []*FeedConfig
}

func NewConfig() *Config {
	c := &Config{
		Verbose:    true,
		Database:   ".feedloggr2.db",
		OutputPath: "feeds",
		Feeds: []*FeedConfig{
			&FeedConfig{"Title of feed", "https://example.com/rss"},
		},
	}

	return c
}

func LoadConfig(path string) (*Config, error) {
	data, e := ioutil.ReadFile(path)
	if e != nil {
		return nil, e
	}
	c := &Config{}
	e = json.Unmarshal(data, &c)
	if e != nil {
		return nil, e
	}
	return c, nil
}

func SaveConfig(path string, c *Config) error {
	b, e := json.MarshalIndent(c, "", "    ")
	if e != nil {
		return e
	}

	e = ioutil.WriteFile(path, b, 0644)
	if e != nil {
		return e
	}
	return nil
}
