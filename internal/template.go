package internal

import (
	"bytes"
	_ "embed"
	"errors"
	html "html/template"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

//go:embed default.html
var defaultTemplate string

const (
	filePerm fs.FileMode = 0644
	dirPerm  fs.FileMode = 0755
)

// TemplateFuncs contains some simple helper functions available inside a template.
var TemplateFuncs = html.FuncMap{
	"shortdate": func(t time.Time) string {
		return t.Format("2006-01-02")
	},
	"prevday": func(t time.Time) time.Time {
		return t.AddDate(0, 0, -1)
	},
	"nextday": func(t time.Time) time.Time {
		return t.AddDate(0, 0, 1)
	},
}

// TemplateGenerator contains the basic info for this generator.
type TemplateGenerator struct {
	Name    string
	Version string
	Source  string
}

// TemplateFeed contains a feed and it's parsed output (items or an error).
type TemplateFeed struct {
	Conf  Feed   // Basic config for the feed
	Items []Item // Any parsed and filtered items
	Error error  // Error returned when trying to download/parse the feed
}

// TemplateVars is a set of basic info that can be provided when executing/writing a template.
type TemplateVars struct {
	Today     time.Time         // Current time
	Generator TemplateGenerator // Basic generator info
	Feeds     []TemplateFeed    // List of feeds and their config, items and errors
}

// NewTemplateVars creates a new instance and adds current time and generator info to it.
// The Feeds field will be empty and has to be manually updated.
func NewTemplateVars() TemplateVars {
	return TemplateVars{
		Today: time.Now(),
		Generator: TemplateGenerator{
			Name:    GeneratorName,
			Version: GeneratorVersion,
			Source:  GeneratorSource,
		},
	}
}

// LoadTemplates tries to parse a template from file or use a default template.
// The returned template has no name and has some helper functions declared.
func LoadTemplate(file string) (tmpl *html.Template, err error) {
	tmpl = html.New("").Funcs(TemplateFuncs)
	if len(file) == 0 {
		tmpl, err = tmpl.Parse(defaultTemplate)
	} else {
		var b []byte
		// gosec warns about file inclusion by variable something, which we kinda wanna do here
		b, err = os.ReadFile(file) // #nosec G304
		if err == nil {
			tmpl, err = tmpl.Parse(string(b))
		}
	}
	return
}

// WriteTemplate executes a loaded template (using provided vars) and writes
// it's output to a file.
func WriteTemplate(file string, tmpl *html.Template, vars interface{}) error {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return err
	}
	d := filepath.Dir(file)
	if err := os.MkdirAll(d, dirPerm); err != nil {
		return err
	}
	return os.WriteFile(file, bytes.TrimSpace(buf.Bytes()), filePerm)
}

// Symlink tries to make a new symlink dst pointing to file src.
func Symlink(src, dst string) error {
	if err := os.Remove(dst); err != nil {
		// Ignore error if the symlink simply doesn't exist yet
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return os.Symlink(filepath.Base(src), dst)
}
