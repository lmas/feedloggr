package internal

import (
	"fmt"
	"io"
	"net/url"
	"regexp"

	"github.com/mmcdole/gofeed"
)

// Basic info about this generator
const (
	GeneratorName    string = "feedloggr"
	GeneratorVersion string = "v0.3.0"
	GeneratorSource  string = "https://github.com/lmas/feedloggr"
)

// Generator contains the runtime state for downloading/parsing/filtering and finally writing news feeds
type Generator struct {
	conf       Conf
	client     *client
	feedParser *gofeed.Parser
	filter     *filter
}

// New creates a new Generator instance, based on conf
func NewGenerator(conf Conf) (*Generator, error) {
	g := &Generator{
		conf: conf,
		client: newClient(clientConf{
			Timeout: conf.Settings.Timeout,
			Jitter:  conf.Settings.Jitter,
		}),
		feedParser: gofeed.NewParser(),
	}
	var err error
	if g.filter, err = loadFilter(conf.Settings.Output); err != nil {
		return nil, err
	}
	return g, nil
}

// NewItems is a shortcut to download/parse/filter a news feed
func (g *Generator) NewItems(f Feed) ([]Item, error) {
	body, err := g.client.RateLimitedGet(f.Url)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var items []Item
	if f.Parser.Rule == "" {
		items, err = g.ParseFeed(body)
	} else {
		items, err = g.ParsePage(body, f)
	}
	if err != nil {
		return nil, err
	}

	filtered := g.filter.filterItems(g.conf.Settings.MaxItems, items...)
	if err = g.filter.write(); err != nil {
		return nil, err
	}
	return filtered, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func newItem(title, url, content string) Item {
	if title == "" {
		if content != "" {
			title = content
		} else {
			title = url // Last ditch thing
		}
	}
	return Item{
		Title:   title,
		Url:     url,
		Content: content,
	}
}

// ParseFeed tries to parse a normal atom/rss/json feed and return it's items
func (g *Generator) ParseFeed(body io.ReadCloser) ([]Item, error) {
	f, err := g.feedParser.Parse(body)
	if err != nil {
		return nil, err
	}
	var items []Item
	for _, i := range f.Items {
		if i.Link == "" {
			continue // You never know, I have low expectations of site feeds...
		}
		if i.Content == "" {
			i.Content = i.Description
		}
		items = append(items, newItem(i.Title, i.Link, i.Content))
	}
	return items, nil
}

// ParsePage sets up a bunch of regexp rules and urls and tries to parse a raw page body for custom items
func (g *Generator) ParsePage(body io.ReadCloser, feed Feed) (items []Item, err error) {
	re, err := regexp.Compile(feed.Parser.Rule)
	if err != nil {
		return
	}
	it, iu, ic := re.SubexpIndex("title"), re.SubexpIndex("url"), re.SubexpIndex("content")
	if it == -1 || iu == -1 {
		err = fmt.Errorf("filter rule missing title or url capture group: %s", feed.Parser.Rule)
		return
	}
	// blanket trust in the Host field, no matter what it's set as (or not set at all)
	feedUrl, err := url.Parse(feed.Parser.Host)
	if err != nil {
		return
	}
	b, err := io.ReadAll(body)
	if err != nil {
		return
	}

	for _, item := range re.FindAllSubmatch(b, -1) {
		var u *url.URL
		u, err = url.Parse(string(item[iu]))
		if err != nil {
			return nil, err
		}
		if u.Scheme == "" || u.Host == "" {
			u = feedUrl.ResolveReference(u)
		}
		var t string
		if it != -1 {
			t = string(item[it])
		}
		var c string
		if ic != -1 {
			c = string(item[ic])
		}
		items = append(items, newItem(t, u.String(), c))
	}

	if len(items) < 1 {
		err = fmt.Errorf("filter rule failed to match any items: %s", feed.Parser.Rule)
	}
	return
}

// FilterStats returns a FilterStats struct with the current state of the internal bloom filter.
func (g *Generator) FilterStats() FilterStats {
	return g.filter.stats()
}
