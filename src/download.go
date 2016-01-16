package feedloggr2

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lmas/go-pkg-xmlx"
)

const USER_AGENT = "feedloggr2/" + VERSION

func parse_feed(url string) ([]*FeedItem, error) {
	data, e := download_feed(url)
	if e != nil {
		return nil, e
	}

	doc := xmlx.New()
	e = doc.LoadString(data, nil)
	if e != nil {
		return nil, e
	}

	var items []*FeedItem
	// TODO: need to come up with a better way to determine the feed type.
	if node := doc.SelectNode("http://www.w3.org/2005/Atom", "feed"); node != nil {
		items, e = parse_atom(url, data)
		if e != nil {
			return nil, e
		}
	} else if node := doc.SelectNode("", "rss"); node != nil {
		// TODO: sometimes the rss tag is not set, but instead some RDF tag?
		items, e = parse_rss(url, data)
		if e != nil {
			return nil, e
		}
	} else {
		return nil, fmt.Errorf("Can't parse feed of unknown type")
	}
	return items, nil
}

func download_feed(url string) (string, error) {
	req, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return "", e
	}

	req.Header.Set("User-Agent", USER_AGENT)
	client := http.DefaultClient
	res, e := client.Do(req)
	if e != nil {
		return "", e
	}

	defer res.Body.Close()
	data, e := ioutil.ReadAll(res.Body)
	if e != nil {
		return "", e
	}

	return string(data), nil
}
