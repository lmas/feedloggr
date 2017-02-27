package feedloggr2

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/mmcdole/gofeed"
)

const (
	UserAgent   string = "feedloggr2/0.2"
	maxItems    int    = 50
	feedTimeout int    = 2 // seconds
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
		Timeout: time.Duration(config.DownloadTimeout) * time.Second,
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
	feeds := app.updateAllFeeds(app.Config.Feeds)

	b, err := app.generatePage(feeds)
	if err != nil {
		return err
	}

	path := filepath.Join(app.Config.OutputPath, date(app.time)+".html")
	err = app.writePage(path, b)
	if err != nil {
		return err
	}

	err = app.writeStyleFile()
	if err != nil {
		return err
	}
	return nil
}
