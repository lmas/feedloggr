package feedloggr

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	filterSize  uint = 1000000 // Should be enough for a couple of years
	maxItems    int  = 50
	feedTimeout int  = 2 // seconds
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func date(t time.Time) string {
	return t.Format("2006-01-02")
}

func (app *App) seenItem(url string) bool {
	return app.filter.Lookup([]byte(url))
}

func (app *App) newItems(url string) ([]Item, error) {
	feed, err := app.feedParser.ParseURL(url)
	if err != nil {
		return nil, err
	}

	var items []Item
	num := min(len(feed.Items), maxItems)
	for _, i := range feed.Items[:num] {
		if app.seenItem(i.Link) {
			continue
		}

		items = append(items, Item{
			Title: strings.TrimSpace(i.Title),
			URL:   i.Link,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Title < items[j].Title
	})
	return items, nil
}

////////////////////////////////////////////////////////////////////////////////

func (app *App) updateAllFeeds(feeds map[string]string) []Feed {
	var updated []Feed
	sleep := time.Duration(feedTimeout) * time.Second
	for title, url := range feeds {
		app.Log("Updating %s (%s)", title, url)
		items, err := app.newItems(url)
		if err != nil {
			app.Log("%s", err)
		}

		if len(items) > 0 || err != nil {
			updated = append(updated, Feed{
				Title: title,
				URL:   url,
				Items: items,
				Error: err,
			})
		}
		time.Sleep(sleep)
	}
	sort.Slice(updated, func(i, j int) bool {
		return updated[i].Title < updated[j].Title
	})
	return updated
}

func (app *App) generatePage(feeds []Feed) ([]byte, error) {
	app.Log("Generating page...")
	var buf bytes.Buffer
	err := app.tmpl.Execute(&buf, map[string]interface{}{
		"CurrentDate": date(app.time),
		"PrevDate":    date(app.time.AddDate(0, 0, -1)),
		"NextDate":    date(app.time.AddDate(0, 0, 1)),
		"Feeds":       feeds,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (app *App) writePage(index, path string, b []byte) error {
	app.Log("Writing page to %s...", path)
	err := os.MkdirAll(filepath.Dir(path), 0744)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return err
	}

	err = os.Remove(index)
	if err != nil {
		// ignore error if the symlink doesn't exist already
		if !os.IsNotExist(err) {
			return err
		}
	}

	err = os.Symlink(filepath.Base(path), index)
	return err
}

func (app *App) writeFilter(path string, feeds []Feed) error {
	for _, f := range feeds {
		for _, i := range f.Items {
			app.filter.Insert([]byte(i.URL))
		}
	}

	err := ioutil.WriteFile(path, app.filter.Encode(), 0644)
	return err
}

func (app *App) writeStyleFile(path string) error {
	// With these flags we try to avoid overwriting an existing file
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}

	defer f.Close()
	app.Log("Writing style file...")
	_, err = f.WriteString(tmplCSS)
	return err
}
