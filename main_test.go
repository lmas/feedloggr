package feedloggr2

import (
	"html/template"
	"net/http"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
)

func makeTestApp(t *testing.T) *App {
	conf := &Config{
		Verbose:         false,
		Database:        ":memory:",
		OutputPath:      "",
		DownloadTimeout: 1,
		Feeds:           []Item{},
	}

	app, err := NewApp(conf)
	if err != nil {
		t.Fatalf("Failed to make new app: %s", err)
	}
	//app.httpClient.Timeout = time.Duration(timeout) * time.Second

	return app
}

func TestNewApp(t *testing.T) {
	makeTestApp(t)
}

func TestApp_Log(t *testing.T) {
	type fields struct {
		Config     *Config
		parser     *gofeed.Parser
		httpClient *http.Client
		db         *DB
		tmpl       *template.Template
		time       time.Time
	}
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				Config:     tt.fields.Config,
				parser:     tt.fields.parser,
				httpClient: tt.fields.httpClient,
				db:         tt.fields.db,
				tmpl:       tt.fields.tmpl,
				time:       tt.fields.time,
			}
			app.Log(tt.args.msg, tt.args.args...)
		})
	}
}

func TestApp_Update(t *testing.T) {
	type fields struct {
		Config     *Config
		parser     *gofeed.Parser
		httpClient *http.Client
		db         *DB
		tmpl       *template.Template
		time       time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				Config:     tt.fields.Config,
				parser:     tt.fields.parser,
				httpClient: tt.fields.httpClient,
				db:         tt.fields.db,
				tmpl:       tt.fields.tmpl,
				time:       tt.fields.time,
			}
			if err := app.Update(); (err != nil) != tt.wantErr {
				t.Errorf("App.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
