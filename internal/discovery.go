package internal

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var discoverNames = []string{
	// These strings have been manually confirmed to be existing in the wild
	// and catches a clear majority of all valid feeds (on a single page).
	"atom",
	"feed",
	"rss",

	// Following strings are hypothetical but haven't been confirmed to catch
	// any feeds that the above strings didn't already find.
	// They also causes too many false positives and are disabled for now.
	// "news",
	// "xml",
}

func inList(item string, list []string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

// DiscoverFeeds tries to discover any URLs that looks like feeds, from a site.
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
