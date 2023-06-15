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

// TmplFuncs contains some custom funcs for being used in the templates
var TmplFuncs = html.FuncMap{
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

type TemplateGenerator struct {
	Name    string
	Version string
	Source  string
}

type TemplateFeed struct {
	Conf  Feed
	Items []Item
	Error error
}

type TemplateVars struct {
	Today     time.Time
	Generator TemplateGenerator
	Feeds     []Feed
}

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

// LoadTemplates returns a html.Template struct, loaded with the parsed templates and ready for use
func LoadTemplate(file string) (*html.Template, error) {
	var err error
	tmpl := html.New("").Funcs(TmplFuncs)
	if len(file) == 0 {
		tmpl, err = tmpl.Parse(defaultTemplate)
	} else {
		var b []byte
		b, err = os.ReadFile(file)
		if err == nil {
			tmpl, err = tmpl.Parse(string(b))
		}
	}
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// WriteTemplate executes a loaded template and writes it's output to a file
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

// Symlink tries to make a new symlink dst pointing to file src
func Symlink(src, dst string) error {
	if err := os.Remove(dst); err != nil {
		// Ignore error if the symlink simply doesn't exist yet
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return os.Symlink(filepath.Base(src), dst)
}
