package feedloggr2

import (
	"fmt"

	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/jteeuwen/go-pkg-xmlx"
)

type FeedDownloader struct {
	rss   *rss.Feed
	items []*FeedItem
	fetch func(string, xmlx.CharsetFunc) error
}

func NewDownloader(user_agent string) *FeedDownloader {
	// Work around the funky model of go-pkg-rss and make a simpler interface.
	f := &FeedDownloader{}
	f.rss = rss.NewWithHandlers(5, false, f, f)
	f.rss.SetUserAgent(user_agent)
	f.fetch = f.rss.Fetch
	return f
}

// Dummy "download" func for testing purposes
// TODO: Move to test file?
func (f *FeedDownloader) dummy_fetch(uri string, charset xmlx.CharsetFunc) error {
	f.rss.Url = uri
	items := []*rss.Item{
		&rss.Item{
			Title: "1st item",
			Links: []*rss.Link{&rss.Link{Href: "http://some.link/"}},
		},
		&rss.Item{
			Title: "2st item",
			Links: []*rss.Link{&rss.Link{Href: "http://some.link2/"}},
		},
	}
	f.ProcessItems(nil, nil, items)
	return nil
}

func (f *FeedDownloader) Clear() {
	f.items = nil
}

func (f *FeedDownloader) DownloadFeed(url string) ([]*FeedItem, error) {
	// TODO: don't actually download a feed when running tests
	defer f.Clear()
	e := f.fetch(url, nil)
	if e != nil {
		return nil, fmt.Errorf("Error connecting to %s: %s\n", url, e)
	}
	return f.items, nil
}

// Dummy func so go-pkg-rss will run.
func (f *FeedDownloader) ProcessChannels(feed *rss.Feed, channels []*rss.Channel) {
}

func (f *FeedDownloader) ProcessItems(feed *rss.Feed, ch *rss.Channel, items []*rss.Item) {
	for _, it := range items {
		f.items = append(f.items, &FeedItem{
			Title: it.Title,
			Url:   it.Links[0].Href,
			Date:  Now(),
			Feed:  f.rss.Url,
		})
	}
}
