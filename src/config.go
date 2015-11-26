package feedloggr2

import (
	"encoding/json"
	"io/ioutil"
)

type FeedConfig struct {
	Title string
	Url   string
}

type Config struct {
	Verbose    bool
	Database   string
	OutputPath string
	Feeds      []FeedConfig
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
