package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Item represents a single news item in a feed.
type Item struct {
	Title   string
	Url     string
	Content string // optional
}

// Parser contains a custom regexp rule for parsing non- atom/rss/json feeds.
type Parser struct {
	Rule string // Regexp rule for gragging items' title/url fields in a feed body
	Host string // Optional prefix the item urls with this host
}

// Feed represents a single news feed and how to download and parse it.
type Feed struct {
	Title  string // Custom title
	Url    string // URL to feed
	Parser Parser `yaml:",omitempty"` // Custom parsing rule.
}

// Source returns the "correct" URL host used as the source for the feed.
func (f Feed) Source() string {
	if f.Parser.Host != "" {
		return f.Parser.Host
	}
	return f.Url
}

// Settings contains the general Generator settings.
type Settings struct {
	Output   string // Dir to output the feeds and internal bloom filter
	Template string // Filepath to custom HTML template
	MaxDays  int    // Max amount of days to keep generated pages for
	MaxItems int    // Max amount of items per feed and per day
	Timeout  int    // Max time in seconds when trying to download a feed
	Jitter   int    // Time in seconds used for randomising rate limits.
	Verbose  bool   // Verbose, debug output
}

// Conf contains ALL settings for a Generator, usually loaded from a yaml file.
type Conf struct {
	Settings Settings // General settings
	Feeds    []Feed   // Per feed settings
}

// LoadConf tries to load a Conf from a yaml file.
func LoadConf(path string) (c Conf, err error) {
	var b []byte
	// gosec warns about file inclusion by variable something, which we kinda wanna do here
	b, err = os.ReadFile(path) // #nosec G304
	if err != nil {
		return
	}

	err = yaml.Unmarshal(b, &c)
	return
}

// ExampleConf returns a working example Conf.
func ExampleConf() Conf {
	return Conf{
		Settings: Settings{
			Output:   "./feeds/",
			Template: "",
			MaxDays:  30,
			MaxItems: 20,
			Timeout:  30,
			Jitter:   2,
			Verbose:  true,
		},
		Feeds: []Feed{
			{
				Title: "Lemmy.link",
				Url:   "https://lemmy.link/feeds/local.xml?sort=TopDay",
			},
			{
				Title: "Hacker News",
				Url:   "https://news.ycombinator.com/rss",
				Parser: Parser{
					Rule: `(?sU)<item>.*<title>(?P<title>[^<]+)</title>.*<comments>(?P<url>[^<]+)</comments>.*</item>`,
					Host: "https://news.ycombinator.com/rss",
				},
			},
		},
	}
}

// String tries to return a yaml formatted string of Conf or an error
func (conf Conf) String() string {
	b, err := yaml.Marshal(conf)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
