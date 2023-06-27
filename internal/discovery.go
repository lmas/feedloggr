package internal

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var discoverNames = []string{
	"rss",
	"atom",
	"feed",
	"news",
	// "xml", // Causes too many false positives, but it may catch something valid..?
}

func inList(item string, list []string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func DiscoverFeeds(site string) (feeds []string, err error) {
	parsed, err := url.Parse(site)
	if err != nil {
		return
	}

	c := newClient(clientConf{
		Timeout: 5,
		Jitter:  0,
	})
	body, err := c.Get(site)
	if err != nil {
		return
	}
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return
	}

	// Loops through all html elements that has an HREF attribute
	doc.Find("*[href]").Each(func(i int, s *goquery.Selection) {
		u, err := url.Parse(s.AttrOr("href", ""))
		if err != nil {
			return
		}
		if u.Scheme == "" {
			u.Scheme = parsed.Scheme
		}
		if u.Host == "" {
			u.Host = parsed.Host
		}

		// Use RequestURI so we can check both the path and query parts of the
		// URL, for example: /bla/bla/?mode=rss
		p := strings.ToLower(u.RequestURI())
		f := strings.ToLower(u.String())
		for _, n := range discoverNames {
			if strings.Contains(p, n) && !inList(f, feeds) {
				feeds = append(feeds, f)
				return
			}
		}
	})
	return
}
