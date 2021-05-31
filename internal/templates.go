package internal

import (
	"bytes"
	"embed"
	"errors"
	html "html/template"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

var (
	//go:embed templates/*
	content embed.FS
	dir     = "templates"
)

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

// LoadTemplates returns a html.Template struct, loaded with the parsed templates and ready for use
func LoadTemplates() (*html.Template, error) {
	tmpls := html.New("")
	err := fs.WalkDir(content, dir, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if de.IsDir() {
			return nil
		}
		b, err := content.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = tmpls.New(filepath.Base(path)).Funcs(TmplFuncs).Parse(string(b))
		return err
	})
	if err != nil {
		return nil, err
	}
	return tmpls, nil
}

// WriteTemplate executes a loaded template and writes it's output to a file
func WriteTemplate(file, name string, tmpls *html.Template, vars interface{}) error {
	var buf bytes.Buffer
	if err := tmpls.ExecuteTemplate(&buf, name, vars); err != nil {
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
