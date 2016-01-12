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

// TODO: set another global version const and use it here
const USER_AGENT = "feedloggr2/0.1"

func UpdateFeeds(c *Config) error {
	db, e := OpenSqliteDB(c.Database)
	if e != nil {
		return e
	}

	if c.Verbose {
		fmt.Println("Downloading feeds...")
	}
	d := NewDownloader(USER_AGENT)
	var all_items []*FeedItem
	for _, f := range c.Feeds {
		items, e := d.DownloadFeed(f.Url)
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

	if c.Verbose {
		fmt.Println("Saving feeds...")
	}
	db.SaveItems(all_items)

	// Iterate over the feeds a 2nd time and grab all saved items for today
	var feeds []*Feed
	for _, f := range c.Feeds {
		items := db.GetItems(f.Url)
		feeds = append(feeds, &Feed{
			Title: f.Title,
			Url:   f.Url,
			Items: items,
		})

		if c.Verbose {
			fmt.Printf("%d items for: %s\n", len(items), f.Title)
		}
	}
	// TODO: must sort the feeds after their names

	if c.Verbose {
		fmt.Println("Generating page...")
	}
	funcmap := template.FuncMap{
		"date_link": func(h int, t time.Time) string {
			d := t.Add(time.Hour * time.Duration(h)).Format("2006-01-02")
			return fmt.Sprintf("%s.html", d)
		},
		"format": func(t time.Time) string { // TODO: rename "format"
			return t.Format("2006-01-02")
		},
	}
	t := template.Must(template.New("TemplateName").Funcs(funcmap).Parse(HTML_BODY))
	s := struct {
		Date  time.Time
		Feeds []*Feed
	}{
		Date:  Now(),
		Feeds: feeds,
	}
	file := fmt.Sprintf("%s.html", Now().Format("2006-01-02"))
	path := filepath.Join(c.OutputPath, file)
	f, e := os.Create(path)
	if e != nil {
		panic(e) // TODO
	}
	e = t.Execute(f, s)
	if e != nil {
		panic(e) // TODO
	}

	if c.Verbose {
		fmt.Println("Updating symlink...")
	}
	path = filepath.Join(c.OutputPath, "index.html")
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

	path = filepath.Join(c.OutputPath, "style.css")
	// With these flags we avoid overwriting an existing file
	f, e = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	defer f.Close()
	if e == nil {
		if c.Verbose {
			fmt.Println("Generating style...")
		}
		_, e = f.WriteString(CSS_BODY)
		if e != nil {
			panic(e) // TODO
		}
	}

	if c.Verbose {
		fmt.Println("Done.")
	}
	return nil
}
