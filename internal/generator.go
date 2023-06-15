package internal

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/mmcdole/gofeed"
)

// Basic info about this generator
const (
	GeneratorName    string = "feedloggr"
	GeneratorVersion string = "v0.3.0"
	GeneratorSource  string = "https://github.com/lmas/feedloggr"
)

type transport struct {
	http.RoundTripper
}

func newTransport(dir string) *transport {
	d := http.DefaultTransport
	// The file protocol enables easier testing
	d.(*http.Transport).RegisterProtocol("file", http.NewFileTransport(http.Dir(dir)))
	return &transport{d}
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", "linux:"+GeneratorName+":"+GeneratorVersion+" ("+GeneratorSource+")")
	return t.RoundTripper.RoundTrip(r)
}

// Generator contains the runtime state for downloading/parsing/filtering and finally writing news feeds
type Generator struct {
	conf       Conf
	client     *http.Client
	feedParser *gofeed.Parser
	filter     *filter
}

// New creates a new Generator instance, based on conf
func New(conf Conf) (*Generator, error) {
	g := &Generator{
		conf: conf,
		client: &http.Client{
			Transport: newTransport("."),
			Timeout:   time.Duration(conf.Settings.Timeout) * time.Second,
		},
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
	body, err := g.Download(f)
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

// Download simply tries to download the body of a feed, using a custom http.Transport
func (g *Generator) Download(feed Feed) (io.ReadCloser, error) {
	r, err := g.client.Get(feed.Url)
	if err != nil {
		return nil, err
	}
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return nil, fmt.Errorf("bad response status: %s", r.Status)
	}
	return r.Body, nil
}

// ParseFeed tries to parse a normal atom/rss/json feed and return it's items
func (g *Generator) ParseFeed(body io.ReadCloser) ([]Item, error) {
	f, err := g.feedParser.Parse(body)
	if err != nil {
		return nil, err
	}
	var items []Item
	for _, i := range f.Items {
		c := i.Content
		if c == "" {
			c = i.Description
		}
		items = append(items, Item{
			Title:   i.Title,
			Url:     i.Link,
			Content: c,
		})
	}
	return items, nil
}

// ParsePage sets up a bunch of regexp rules and urls and tries to parse a raw page body for custom items
func (g *Generator) ParsePage(body io.ReadCloser, feed Feed) ([]Item, error) {
	re, err := regexp.Compile(feed.Parser.Rule)
	if err != nil {
		return nil, err
	}
	it, iu, ic := re.SubexpIndex("title"), re.SubexpIndex("url"), re.SubexpIndex("content")
	if it == -1 || iu == -1 {
		return nil, fmt.Errorf("filter rule missing title or url capture group: %s", feed.Parser.Rule)
	}
	// blanket trust in the Host field, no matter what it's set as (or not set at all)
	feedUrl, err := url.Parse(feed.Parser.Host)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var items []Item
	for _, item := range re.FindAllSubmatch(b, -1) {
		u, err := url.Parse(string(item[iu]))
		if err != nil {
			return nil, err
		}
		if u.Scheme == "" || u.Host == "" {
			u = feedUrl.ResolveReference(u)
		}
		var c string
		if ic != -1 {
			c = string(item[ic])
		}
		items = append(items, Item{
			Title:   string(item[it]),
			Url:     u.String(),
			Content: c,
		})
	}
	return items, nil
}

// FilterStats returns a FilterStats struct with the current state of the internal bloom filter.
func (g *Generator) FilterStats() FilterStats {
	return g.filter.stats()
}
