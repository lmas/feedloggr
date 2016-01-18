package feedloggr2

import (
	"encoding/xml"
	"strings"
)

type AtomFeed struct {
	XMLName xml.Name    `xml:"http://www.w3.org/2005/Atom feed"`
	Items   []*AtomItem `xml:"entry"`
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
			URL:   url,
			Date:  Now(),
			Feed:  url,
		})
	}
	return items, nil
}
