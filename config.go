package feedloggr

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Verbose    bool
	OutputPath string
	Timeout    int // In seconds
	Feeds      map[string]string
}

func (c *Config) String() string {
	var b bytes.Buffer
	err := toml.NewEncoder(&b).Encode(c)
	if err != nil {
		return err.Error()
	}
	return b.String()
}

func NewConfig() *Config {
	c := &Config{
		Verbose:    true,
		OutputPath: "./feeds",
		Timeout:    60,
		Feeds: map[string]string{
			"Example": "https://example.com/rss",
		},
	}
	return c
}

func LoadConfig(path string) (*Config, error) {
	c := &Config{}
	_, err := toml.DecodeFile(path, &c)
	return c, err
}
