package feedloggr2

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func (app *App) updateAllFeeds(feeds []Item) []Feed {
	var updated []Feed
	sleep := time.Duration(feedTimeout) * time.Second
	for _, f := range feeds {
		items, err := app.updateSingleFeed(f)
		if err != nil {
			app.Log("Error: %s", err)
		}

		if len(items) > 0 || err != nil {
			updated = append(updated, Feed{
				Title: f.Title,
				URL:   f.URL,
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

func (app *App) updateSingleFeed(feed Item) ([]Item, error) {
	app.Log("Updating %s (%s)", feed.Title, feed.URL)
	b, err := app.downloadFeed(feed.URL)
	if err != nil {
		return nil, err
	}

	items, err := app.parseFeed(b)
	if err != nil {
		return nil, err
	}

	d := date(app.time)
	err = app.db.SaveItems(feed.URL, d, items)
	if err != nil {
		return nil, err
	}

	uniqe, err := app.db.GetItems(feed.URL, d)
	if err != nil {
		return nil, err
	}

	return uniqe, nil
}

func (app *App) downloadFeed(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)
	res, err := app.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func (app *App) parseFeed(b io.ReadCloser) ([]Item, error) {
	defer b.Close()
	feed, err := app.parser.Parse(b)
	if err != nil {
		return nil, err
	}

	var seen []string
	var items []Item
	num := min(len(feed.Items), maxItems)
	for _, i := range feed.Items[:num] {
		// Avoid items with duplicate names (/r/WorldNews ffs)
		if inList(i.Title, seen) {
			continue
		}
		seen = append(seen, i.Title)

		items = append(items, Item{
			Title: i.Title,
			URL:   i.Link,
		})
	}
	return items, nil
}

func (app *App) generatePage(feeds []Feed) ([]byte, error) {
	app.Log("Generating page...")
	buf := new(bytes.Buffer)
	err := app.tmpl.Execute(buf, map[string]interface{}{
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

func (app *App) writePage(path string, b []byte) error {
	app.Log("Writing page to %s...", path)
	err := os.MkdirAll(filepath.Dir(path), 0744)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return err
	}

	index := filepath.Join(app.Config.OutputPath, "index.html")
	err = os.Remove(index)
	if err != nil {
		// ignore error if the symlink doesn't exist already
		if !os.IsNotExist(err) {
			return err
		}
	}

	err = os.Symlink(filepath.Base(path), index)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) writeStyleFile() error {
	path := filepath.Join(app.Config.OutputPath, "style.css")
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
	if err != nil {
		return err
	}
	return nil
}
