package feedloggr2

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
)

func makeTestServer(f func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(f))
}

func TestApp_updateAllFeeds(t *testing.T) {
	type fields struct {
		Config     *Config
		parser     *gofeed.Parser
		httpClient *http.Client
		db         *DB
		tmpl       *template.Template
		time       time.Time
	}
	type args struct {
		feeds []Item
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Feed
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
			if got := app.updateAllFeeds(tt.args.feeds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("App.updateAllFeeds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_updateSingleFeed(t *testing.T) {
	type fields struct {
		Config     *Config
		parser     *gofeed.Parser
		httpClient *http.Client
		db         *DB
		tmpl       *template.Template
		time       time.Time
	}
	type args struct {
		feed Item
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Item
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
			got, _, err := app.updateSingleFeed(tt.args.feed)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.updateSingleFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("App.updateSingleFeed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_downloadFeed(t *testing.T) {
	app := makeTestApp(t)
	s := makeTestServer(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/good":
			fmt.Fprintf(w, "good response")
		case "/empty":
			fmt.Fprintf(w, "")
		case "/timeout":
			time.Sleep(time.Duration(app.Config.DownloadTimeout+1) * time.Second)
			fmt.Fprintf(w, "timeout response?")
		default:
			fmt.Fprintf(w, "404: Not Found")
		}
	})
	defer s.Close()

	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{"good", s.URL + "/good", "good response", false},
		{"empty response", s.URL + "/empty", "", false},
		{"empty url", "", "", true},
		{"timeout", s.URL + "/timeout", "timeout response?", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _, err := app.downloadFeed(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.downloadFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && tt.wantErr {
				return
			}
			defer r.Close()

			b, err := ioutil.ReadAll(r)
			if err != nil {
				t.Errorf("ioutil.ReadAll(App.downloadFeed()) error = %v", err)
			}
			got := string(b)
			if got != tt.want {
				t.Errorf("App.downloadFeed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_parseFeed(t *testing.T) {
	app := makeTestApp(t)
	tests := []struct {
		file    string
		wantErr bool
	}{
		{"reddit", false},
		{"gp", false},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			testfile := filepath.Join("testdata", tt.file+".feed")
			wantfile := filepath.Join("testdata", tt.file+".wanted")
			tf, err := os.Open(testfile)
			if err != nil {
				t.Fatalf("os.Open(%s) error = %v", testfile, err)
			}
			b, err := ioutil.ReadFile(wantfile)
			if err != nil {
				t.Fatalf("ioutil.ReadFile(%s) error = %v", wantfile, err)
			}
			var wanted []Item
			err = json.Unmarshal(b, &wanted)
			if err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			got, err := app.parseFeed(tf)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.parseFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, wanted) {
				t.Errorf("App.parseFeed() = \n%v\n\nwant\n\n%v", got, wanted)
			}
		})
	}
}

func TestApp_generatePage(t *testing.T) {
	type fields struct {
		Config     *Config
		parser     *gofeed.Parser
		httpClient *http.Client
		db         *DB
		tmpl       *template.Template
		time       time.Time
	}
	type args struct {
		feeds []Feed
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
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
			got, err := app.generatePage(tt.args.feeds)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.generatePage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("App.generatePage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_writePage(t *testing.T) {
	type fields struct {
		Config     *Config
		parser     *gofeed.Parser
		httpClient *http.Client
		db         *DB
		tmpl       *template.Template
		time       time.Time
	}
	type args struct {
		path string
		b    []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
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
			if err := app.writePage(tt.args.path, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("App.writePage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_writeStyleFile(t *testing.T) {
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
			if err := app.writeStyleFile(); (err != nil) != tt.wantErr {
				t.Errorf("App.writeStyleFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
