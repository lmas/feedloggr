package feedloggr2

import (
	"fmt"

	rss "github.com/jteeuwen/go-pkg-rss"
)

type FeedDownloader struct {
	rss   *rss.Feed
	items []*FeedItem
}

func NewDownloader(user_agent string) *FeedDownloader {
	// Work around the funky model of go-pkg-rss and make a simpler interface.
	f := &FeedDownloader{}
	f.rss = rss.NewWithHandlers(5, false, f, f)
	f.rss.SetUserAgent(user_agent)
	return f
}

func (f *FeedDownloader) Clear() {
	f.items = nil
}

func (f *FeedDownloader) DownloadFeed(url string) ([]*FeedItem, error) {
	// TODO: don't actually download a feed when running tests
	defer f.Clear()
	e := f.rss.Fetch(url, nil)
	if e != nil {
		return nil, fmt.Errorf("Error connecting to %s: %s\n", url, e)
	}
	return f.items, nil
}

func (f *FeedDownloader) ProcessChannels(feed *rss.Feed, channels []*rss.Channel) {
	// Dummy func so go-pkg-rss will run.
}

func (f *FeedDownloader) ProcessItems(feed *rss.Feed, ch *rss.Channel, items []*rss.Item) {
	for _, it := range items {
		f.items = append(f.items, &FeedItem{
			Title: it.Title,
			Url:   it.Links[0].Href,
			Date:  Now(),
			Feed:  f.rss.Url,
			//Feed:  feed.Url, // TODO: fix the feed url thing
		})
	}
}
