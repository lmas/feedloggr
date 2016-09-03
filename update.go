package feedloggr2

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"html/template"
)

// TODO: change to env var/flag instead
const timeBetweenFeeds int = 2 // In seconds

func UpdateFeeds(c *Config) error {
	db, e := OpenSqliteDB(c.Database)
	if e != nil {
		return e
	}

	u := &UpdateInstance{
		Config:    c,
		DB:        db,
		bad_feeds: make(map[*FeedConfig]error),
	}

	return u.run()
}

type UpdateInstance struct {
	Config *Config
	DB     *DB

	bad_feeds map[*FeedConfig]error
}

func (u *UpdateInstance) add_bad_feed(f *FeedConfig, e error) {
	u.bad_feeds[f] = e
}

func (u *UpdateInstance) log(s string, args ...interface{}) {
	if u.Config.Verbose {
		fmt.Printf(s+"\n", args...)
	}
}

func (u *UpdateInstance) run() error {
	u.download_feeds()
	feeds := u.get_feeds()
	e := u.generate_page(feeds)
	if e != nil {
		return e
	}
	e = u.generate_style()
	if e != nil {
		return e
	}
	u.log("Done.")
	return nil
}

func (u *UpdateInstance) download_feeds() {
	u.log("Downloading feeds...")
	var all_items []*FeedItem
	for _, f := range u.Config.Feeds {
		items, e := parse_feed(f.URL)
		if e != nil {
			u.add_bad_feed(f, e)
			continue
		}

		if len(items) < 1 {
			continue
		}
		all_items = append(all_items, items...)

		// Slow down the amount of requests, to ensure we won't get spam blocked.
		time.Sleep(time.Duration(timeBetweenFeeds) * time.Second)
	}
	u.log("Saving feeds...")
	u.DB.SaveItems(all_items)
}

func (u *UpdateInstance) get_feeds() []*Feed {
	// Iterate over the feeds a 2nd time and grab all saved items for today
	u.log("Getting today's news...")
	var all_feeds FeedSlice
	for _, f := range u.Config.Feeds {
		if e, ok := u.bad_feeds[f]; ok == true {
			all_feeds = append(all_feeds, &Feed{
				Title: f.Title,
				URL:   f.URL,
				Error: e,
			})
			continue
		}

		items := u.DB.GetItems(f.URL)
		if len(items) < 1 {
			continue
		}

		all_feeds = append(all_feeds, &Feed{
			Title: f.Title,
			URL:   f.URL,
			Items: items,
		})

		u.log("%d items for: %s", len(items), f.Title)
	}
	sort.Sort(all_feeds)
	return all_feeds
}

func (u *UpdateInstance) generate_page(feeds []*Feed) error {
	u.log("Generating page...")
	funcmap := template.FuncMap{
		"date_link": func(h int, t time.Time) string {
			d := t.Add(time.Hour * time.Duration(h)).Format("2006-01-02")
			return fmt.Sprintf("%s.html", d)
		},
		"pretty_date": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
	}

	t := template.Must(template.New("Page").Funcs(funcmap).Parse(htmlBody))
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
		return e
	}
	e = t.Execute(f, s)
	if e != nil {
		return e
	}

	u.log("Updating symlink...")
	path = filepath.Join(u.Config.OutputPath, "index.html")
	e = os.Remove(path)
	if e != nil {
		perr, ok := e.(*os.PathError)
		// Ignore any "no such file" errors
		// It works correctly, but "|| perr.Err == "is logically wrong. Bug?
		if ok == false || perr.Err == fmt.Errorf("no such file or directory") {
			return e
		}
	}
	e = os.Symlink(file, path)
	if e != nil {
		return e
	}

	return nil
}

func (u *UpdateInstance) generate_style() error {
	path := filepath.Join(u.Config.OutputPath, "style.css")
	// With these flags we avoid overwriting an existing file
	f, e := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	defer f.Close()
	if e == nil {
		u.log("Generating style...")
		_, e = f.WriteString(cssBody)
		if e != nil {
			return e
		}
	}

	return nil
}
