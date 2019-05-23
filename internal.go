package feedloggr

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tdl "github.com/lmas/Damerau-Levenshtein"
	boom "github.com/tylertreat/BoomFilters"
)

const (
	maxItems    int = 50
	feedTimeout int = 2 // seconds
)

func date(t time.Time) string {
	return t.Format("2006-01-02")
}

func loadFilter(path string) (*boom.ScalableBloomFilter, error) {
	filter := boom.NewDefaultScalableBloomFilter(0.01)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return filter, nil
		}
		return nil, err
	}

	defer f.Close()
	if _, err := filter.ReadFrom(f); err != nil {
		return nil, err
	}
	return filter, nil
}

////////////////////////////////////////////////////////////////////////////////

func seenTitle(title string, list []string) bool {
	for _, s := range list {
		score := tdl.Distance(title, s)
		if score < 2 {
			return true
		}
	}
	return false
}

func (app *App) seenURL(url string) bool {
	return app.filter.TestAndAdd([]byte(url))
}

func (app *App) newItems(url string) ([]Item, error) {
	feed, err := app.feedParser.ParseURL(url)
	if err != nil {
		return nil, err
	}

	var titles []string
	var items []Item
	max := len(feed.Items)
	if max > maxItems {
		max = maxItems
	}

	for _, i := range feed.Items[:max] {
		if seenTitle(i.Title, titles) {
			continue
		}
		titles = append(titles, i.Title)

		if app.seenURL(i.Link) {
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

func (app *App) generatePage(feeds []Feed) (*bytes.Buffer, error) {
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
	return &buf, nil
}

func (app *App) writePage(index, path string, buf *bytes.Buffer) error {
	app.Log("Writing page to %s...", path)
	err := os.MkdirAll(filepath.Dir(path), 0744)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := buf.WriteTo(f); err != nil {
		return err
	}

	err = os.Remove(index)
	if err != nil {
		// ignore error if the symlink doesn't exist already
		if !os.IsNotExist(err) {
			return err
		}
	}

	return os.Symlink(filepath.Base(path), index)
}

func (app *App) writeFilter(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = app.filter.WriteTo(f)
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
