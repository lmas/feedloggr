package feedloggr2

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const userAgent string = "feedloggr2/" + Version

func parse_feed(url string) ([]*FeedItem, error) {
	data, e := download_feed(url)
	if e != nil {
		return nil, e
	}

	var items []*FeedItem
	// First attempt to parse the feed as RSS
	items, e = parse_rss(url, data)
	if e == nil {
		return items, nil
	}

	// If that fails try to parse it as Atom
	items, e = parse_atom(url, data)
	if e == nil {
		return items, nil
	}

	// Or give up
	return nil, fmt.Errorf("Can't parse feed: %v", e)
}

func download_feed(url string) ([]byte, error) {
	req, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return nil, e
	}

	req.Header.Set("User-Agent", userAgent)
	client := http.DefaultClient
	res, e := client.Do(req)
	if e != nil {
		return nil, e
	}

	defer res.Body.Close()
	data, e := ioutil.ReadAll(res.Body)
	if e != nil {
		return nil, e
	}

	return data, nil
}
