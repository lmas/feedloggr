package feedloggr2

import (
	"encoding/xml"
	"strings"
)

type AtomFeed struct {
	Items []*AtomItem `xml:"entry"`
}

type AtomItem struct {
	Title string      `xml:"title"`
	Links []*AtomLink `xml:"link"`
}

type AtomLink struct {
	URL string `xml:"href,attr"`
}

func parse_atom(url, body string) ([]*FeedItem, error) {
	f := AtomFeed{}
	decoder := xml.NewDecoder(strings.NewReader(body))
	e := decoder.Decode(&f)
	if e != nil {
		return nil, e
	}

	var items []*FeedItem
	for _, i := range f.Items {
		url := i.Links[0].URL
		items = append(items, &FeedItem{
			Title: i.Title,
			Url:   url,
			Date:  Now(),
			Feed:  url,
		})
	}
	return items, nil
}
