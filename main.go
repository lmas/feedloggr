package feedloggr2

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/mmcdole/gofeed"
)

const (
	UserAgent       string = "feedloggr2/0.2"
	maxItems        int    = 50
	downloadTimeout int    = 60 // seconds
	feedTimeout     int    = 2  // seconds
)

// Item is used by Config.Feeds and Feed.Items
type Item struct {
	Title string
	URL   string
}

type Feed struct {
	Title string
	URL   string
	Items []Item
	Error error
}

type App struct {
	Config     *Config
	parser     *gofeed.Parser
	httpClient *http.Client
	db         *DB
	tmpl       *template.Template
	time       time.Time
}

func New(config *Config) (*App, error) {
	db, err := OpenDB(config.Database)
	if err != nil {
		return nil, err
	}

	parser := gofeed.NewParser()
	client := &http.Client{
		Timeout: time.Duration(downloadTimeout) * time.Second,
	}
	tmpl, err := template.New("page").Parse(tmplPage)
	if err != nil {
		return nil, err
	}

	app := &App{
		Config:     config,
		parser:     parser,
		httpClient: client,
		db:         db,
		tmpl:       tmpl,
		time:       time.Now(),
	}
	return app, nil
}

func (app *App) Log(msg string, args ...interface{}) {
	if app.Config.Verbose {
		log.Printf(msg+"\n", args...)
	}
}

func (app *App) Update() error {
	var feeds []Feed
	sleep := time.Duration(feedTimeout) * time.Second
	for _, f := range app.Config.Feeds {
		items, err := app.updateFeed(f)
		if err != nil {
			app.Log("Error: %s", err)
		}
		feeds = append(feeds, Feed{
			Title: f.Title,
			URL:   f.URL,
			Items: items,
			Error: err,
		})
		time.Sleep(sleep)
	}
	sort.Slice(feeds, func(i, j int) bool {
		return feeds[i].Title < feeds[j].Title
	})

	err := app.generatePage(feeds)
	if err != nil {
		return err
	}
	err = app.generateStyle()
	if err != nil {
		return err
	}
	return nil
}

func (app *App) updateFeed(feed Item) ([]Item, error) {
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

func (app *App) generatePage(feeds []Feed) error {
	app.Log("Generating page...")
	today := date(app.time)

	err := os.MkdirAll(app.Config.OutputPath, 0744)
	if err != nil {
		return err
	}

	page := today + ".html"
	path := filepath.Join(app.Config.OutputPath, page)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = app.tmpl.Execute(f, map[string]interface{}{
		"CurrentDate": today,
		"PrevDate":    date(app.time.AddDate(0, 0, -1)),
		"NextDate":    date(app.time.AddDate(0, 0, 1)),
		"Feeds":       feeds,
	})
	if err != nil {
		return err
	}

	app.Log("Updating symlink")
	index := filepath.Join(app.Config.OutputPath, "index.html")
	err = os.Remove(index)
	if err != nil {
		// ignore error if the symlink doesn't exist already
		if !os.IsNotExist(err) {
			return err
		}
	}

	err = os.Symlink(page, index)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) generateStyle() error {
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
	app.Log("Generating style")
	_, err = f.WriteString(tmplCSS)
	if err != nil {
		return err
	}
	return nil
}
