package feedloggr2

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"html/template"
)

// TODO: change to env var/flag instead
const TIME_BETWEEN_FEEDS = 2 // In seconds

func UpdateFeeds(c *Config) error {
	db, e := OpenSqliteDB(c.Database)
	if e != nil {
		return e
	}

	d := NewDownloader("feedloggr2/" + VERSION)

	u := &UpdateInstance{
		Config:     c,
		DB:         db,
		Downloader: d,
	}

	return u.run()
}

type UpdateInstance struct {
	Config     *Config
	DB         *DB
	Downloader *FeedDownloader
}

func (u *UpdateInstance) log(s string, args ...interface{}) {
	if u.Config.Verbose {
		fmt.Printf(s+"\n", args...)
	}
}

func (u *UpdateInstance) run() error {
	u.download_feeds()
	feeds := u.get_feeds()
	u.generate_page(feeds)
	u.generate_style()
	u.log("Done.")
	return nil
}

func (u *UpdateInstance) download_feeds() {
	u.log("Downloading feeds...")
	var all_items []*FeedItem
	for _, f := range u.Config.Feeds {
		items, e := u.Downloader.DownloadFeed(f.Url)
		if e != nil {
			fmt.Println(e)
			continue
		}

		if len(items) < 1 {
			continue
		}
		all_items = append(all_items, items...)

		// Slow down the amount of requests, to ensure we won't get spam blocked.
		time.Sleep(time.Duration(TIME_BETWEEN_FEEDS) * time.Second)
	}
	u.log("Saving feeds...")
	u.DB.SaveItems(all_items)
}

func (u *UpdateInstance) get_feeds() []*Feed {
	// Iterate over the feeds a 2nd time and grab all saved items for today
	u.log("Getting today's news...")
	var all_feeds []*Feed
	for _, f := range u.Config.Feeds {
		items := u.DB.GetItems(f.Url)
		all_feeds = append(all_feeds, &Feed{
			Title: f.Title,
			Url:   f.Url,
			Items: items,
		})

		u.log("%d items for: %s", len(items), f.Title)
	}
	// TODO: must sort the feeds after their names
	return all_feeds
}

func (u *UpdateInstance) generate_page(feeds []*Feed) {
	u.log("Generating page...")
	funcmap := template.FuncMap{
		"date_link": func(h int, t time.Time) string {
			d := t.Add(time.Hour * time.Duration(h)).Format("2006-01-02")
			return fmt.Sprintf("%s.html", d)
		},
		"format": func(t time.Time) string { // TODO: rename "format"
			return t.Format("2006-01-02")
		},
	}

	t := template.Must(template.New("Page").Funcs(funcmap).Parse(HTML_BODY))
	s := struct {
		Date  time.Time
		Feeds []*Feed
	}{
		Date:  Now(),
		Feeds: feeds,
	}
	file := fmt.Sprintf("%s.html", Now().Format("2006-01-02"))
	path := filepath.Join(u.Config.OutputPath, file)
	f, e := os.Create(path)
	defer f.Close()
	if e != nil {
		panic(e) // TODO
	}
	e = t.Execute(f, s)
	if e != nil {
		panic(e) // TODO
	}

	u.log("Updating symlink...")
	path = filepath.Join(u.Config.OutputPath, "index.html")
	e = os.Remove(path)
	if e != nil {
		perr, ok := e.(*os.PathError)
		// Ignore any "no such file" errors
		// It works correctly, but "|| perr.Err == "is logically wrong. Bug?
		if ok == false || perr.Err == fmt.Errorf("no such file or directory") {
			panic(e) // TODO
		}
	}
	e = os.Symlink(file, path)
	if e != nil {
		panic(e) // TODO
	}
}

func (u *UpdateInstance) generate_style() {
	path := filepath.Join(u.Config.OutputPath, "style.css")
	// With these flags we avoid overwriting an existing file
	f, e := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	defer f.Close()
	if e == nil {
		u.log("Generating style...")
		_, e = f.WriteString(CSS_BODY)
		if e != nil {
			panic(e) // TODO
		}
	}
}
