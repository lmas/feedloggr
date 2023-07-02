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
	GeneratorVersion string = "v0.4.0"
	GeneratorSource  string = "https://github.com/lmas/feedloggr"
)

// Generator is used for downloading, parsing and then filtering items from feeds.
type Generator struct {
	conf       Conf
	client     *client
	feedParser *gofeed.Parser
	filter     *filter
}

// New creates a new Generator instance, based on conf.
func NewGenerator(conf Conf) (gen *Generator, err error) {
	gen = &Generator{
		conf: conf,
		client: newClient(clientConf{
			Timeout: conf.Settings.Timeout,
			Jitter:  conf.Settings.Jitter,
		}),
		feedParser: gofeed.NewParser(),
	}
	gen.filter, err = loadFilter(conf.Settings.Output)
	return
}

// FetchItems downloads a feed and tries to find any items in it.
func (g *Generator) FetchItems(f Feed) (items []Item, err error) {
	var body io.ReadCloser
	get := g.client.Get
	if g.conf.Settings.Jitter > 1 {
		get = g.client.RateLimitedGet
	}

	body, err = get(f.Url)
	if err != nil {
		return
	}
	defer body.Close()

	if f.Parser.Rule == "" {
		items, err = g.parseFeed(body)
	} else {
		items, err = g.parseFeedRegexp(body, f.Parser)
	}
	return
}

// NewItems downloads a feed, tries to find any items and filter out the ones
// that has already been seen before.
func (g *Generator) NewItems(f Feed) (items []Item, err error) {
	items, err = g.FetchItems(f)
	if err != nil || len(items) < 1 {
		return
	}
	items = g.filter.filterItems(g.conf.Settings.MaxItems, items...)
	err = g.filter.write()
	return
}

// FilterStats returns a FilterStats struct with the current state of the
// internal bloom filter.
func (g *Generator) FilterStats() FilterStats {
	return g.filter.stats()
}

////////////////////////////////////////////////////////////////////////////////

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

// parseFeed tries to parse a normal atom/rss/json feed and return it's items
func (g *Generator) parseFeed(body io.ReadCloser) (items []Item, err error) {
	f, err := g.feedParser.Parse(body)
	if err != nil {
		return
	}

	for _, i := range f.Items {
		if i.Link == "" {
			continue // You never know, I have low expectations of site feeds...
		}
		if i.Content == "" {
			i.Content = i.Description
		}
		items = append(items, newItem(i.Title, i.Link, i.Content))
	}
	return
}

// parseFeedRegexp compiles a regexp rule and tries to parse a raw page body for custom items
func (g *Generator) parseFeedRegexp(body io.ReadCloser, parser Parser) (items []Item, err error) {
	re, err := regexp.Compile(parser.Rule)
	if err != nil {
		return
	}
	it, iu, ic := re.SubexpIndex("title"), re.SubexpIndex("url"), re.SubexpIndex("content")
	if it == -1 || iu == -1 {
		err = fmt.Errorf("filter rule missing title or url capture group: %s", parser.Rule)
		return
	}
	// blanket trust in the Host field, no matter what it's set as (or not set at all)
	feedUrl, err := url.Parse(parser.Host)
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
		err = fmt.Errorf("filter rule failed to match any items: %s", parser.Rule)
	}
	return
}
