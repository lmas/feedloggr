package main

import "github.com/lmas/feedloggr2/src"

func main() {
	// TODO: make example config
	feeds := []feedloggr2.FeedConfig{
		feedloggr2.FeedConfig{"reddit - Front", "https://reddit.com/.rss"},
		feedloggr2.FeedConfig{"reddit - Funny", "https://reddit.com/r/funny.rss"},
	}
	c := &feedloggr2.Config{
		Verbose:    true,
		Database:   "./tmp/feeds.db",
		OutputPath: "./tmp/out/",
		Feeds:      feeds,
	}
	feedloggr2.SaveConfig("./tmp/config.json", c)

	// Load config
	config, e := feedloggr2.LoadConfig("./tmp/config.json")
	if e != nil {
		panic(e) // TODO
	}

	// Update
	e = feedloggr2.UpdateFeeds(config)
	if e != nil {
		panic(e) // TODO
	}
}
