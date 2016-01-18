package feedloggr2

import (
	"encoding/xml"
	"strings"
)

type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Items []*RSSItem `xml:"item"`
}

type RSSItem struct {
	Title string `xml:"title"`
	URL   string `xml:"link"`
}

func parse_rss(url, body string) ([]*FeedItem, error) {
	f := RSSFeed{}
	decoder := xml.NewDecoder(strings.NewReader(body))
	e := decoder.Decode(&f)
	if e != nil {
		return nil, e
	}

	var items []*FeedItem
	for _, i := range f.Channel.Items {
		items = append(items, &FeedItem{
			Title: i.Title,
			URL:   i.URL,
			Date:  Now(),
			Feed:  url,
		})
	}

	return items, nil
}
