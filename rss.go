package feedloggr2

import "encoding/xml"

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

func parse_rss(url string, body []byte) ([]*FeedItem, error) {
	f := RSSFeed{}
	e := xml.Unmarshal(body, &f)
	if e != nil {
		return nil, e
	}

	var items []*FeedItem
	for _, i := range f.Channel.Items {
		items = append(items, &FeedItem{
			Title: i.Title,
			URL:   i.URL,
			Date:  Now(), // TODO: remove this, added in db.go
			Feed:  url,
		})
	}

	return items, nil
}
