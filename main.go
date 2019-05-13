package feedloggr

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/mmcdole/gofeed"
	cuckoo "github.com/seiflotfy/cuckoofilter"
)

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
	Config *Config

	time       time.Time
	tmpl       *template.Template
	filter     *cuckoo.Filter
	feedParser *gofeed.Parser
}

func New(config *Config) (*App, error) {
	tmpl, err := template.New("page").Parse(tmplPage)
	if err != nil {
		return nil, err
	}

	feedParser := gofeed.NewParser()
	feedParser.Client = &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	app := &App{
		Config:     config,
		time:       time.Now(),
		tmpl:       tmpl,
		feedParser: feedParser,
	}

	path := filepath.Join(config.OutputPath, ".filter.dat")
	b, err := ioutil.ReadFile(path)
	if err == nil {
		app.filter, err = cuckoo.Decode(b)
		if err != nil {
			return nil, err
		}
	} else {
		app.Log("Error loading filter: %s", err)
		app.filter = cuckoo.NewFilter(filterSize)
	}
	return app, nil
}

func (app *App) Log(msg string, args ...interface{}) {
	if app.Config.Verbose {
		log.Printf(msg+"\n", args...)
	}
}

func (app *App) Update() error {
	feeds := app.updateAllFeeds(app.Config.Feeds)
	b, err := app.generatePage(feeds)
	if err != nil {
		return err
	}

	index := filepath.Join(app.Config.OutputPath, "index.html")
	path := filepath.Join(app.Config.OutputPath, date(app.time)+".html")
	if err := app.writePage(index, path, b); err != nil {
		return err
	}

	path = filepath.Join(app.Config.OutputPath, ".filter.dat")
	if err := app.writeFilter(path, feeds); err != nil {
		return err
	}

	path = filepath.Join(app.Config.OutputPath, "style.css")
	if err := app.writeStyleFile(path); err != nil {
		return err
	}
	return nil
}
