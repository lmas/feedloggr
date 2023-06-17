package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/lmas/feedloggr/internal"
)

var (
	confFile    = flag.String("conf", ".feedloggr.yml", "Path to conf file")
	confExample = flag.Bool("example", false, "Print example config and exit")
	confTest    = flag.Bool("test", false, "Load config and exit")
	confVerbose = flag.Bool("verbose", false, "Print debug messages while running")
	confVersion = flag.Bool("version", false, "Print version and exit")
)

func main() {
	flag.Parse()

	// Early quitters
	switch {
	case *confExample:
		fmt.Println(internal.ExampleConf())
		os.Exit(0)
	case *confVersion:
		fmt.Printf("%s %s\n", internal.GeneratorName, internal.GeneratorVersion)
		os.Exit(0)
	}

	conf, err := internal.LoadConf(*confFile)
	if err != nil {
		panic(err)
	}

	// Late quitter
	if *confTest {
		fmt.Println(conf)
		fmt.Println("No errors while loading config")
		os.Exit(0)
	}

	if *confVerbose != conf.Settings.Verbose {
		// Weeell if one of 'em is true dey bath gotta be true nao
		*confVerbose, conf.Settings.Verbose = true, true
	}

	tmpl, err := internal.LoadTemplate(conf.Settings.Template)
	if err != nil {
		panic(err)
	}

	feeds, err := fetchFeeds(conf)
	if err != nil {
		panic(err)
	}

	if err := writeFiles(conf.Settings.Output, feeds, tmpl); err != nil {
		panic(err)
	}

	if err := removeOldFiles(conf.Settings.Output, conf.Settings.MaxDays); err != nil {
		panic(err)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func debug(msg string, args ...interface{}) {
	if *confVerbose {
		fmt.Printf(msg+"\n", args...)
	}
}

func fetchFeeds(conf internal.Conf) (feeds []internal.TemplateFeed, err error) {
	gen, err := internal.NewGenerator(conf)
	if err != nil {
		return
	}

	for _, source := range conf.Feeds {
		debug("Updating %s (%s)", source.Title, source.Url)
		items, err := gen.NewItems(source)
		if err != nil {
			debug("\tError: %s", err)
		} else {
			debug("\tItems: %d", len(items))
		}

		if len(items) > 0 || err != nil {
			feeds = append(feeds, internal.TemplateFeed{
				Conf:  source,
				Items: items,
				Error: err,
			})
		}
	}

	debug("Filter stats: %+v\n", gen.FilterStats())
	return
}

func writeFiles(dir string, feeds []internal.TemplateFeed, tmpl *template.Template) error {
	v := internal.NewTemplateVars()
	v.Feeds = feeds
	p := filepath.Join(dir, "news-"+v.Today.Format("2006-01-02")+".html")
	if err := internal.WriteTemplate(p, tmpl, v); err != nil {
		return err
	}
	debug("Wrote %s", p)
	if err := internal.Symlink(p, filepath.Join(dir, "index.html")); err != nil {
		return err
	}
	return nil
}

var reFile = regexp.MustCompile(`^.*/news-(\d\d\d\d-\d\d-\d\d).html$`)

func removeOldFiles(dir string, maxDays int) error {
	if maxDays < 1 {
		return nil
	}

	cutoff := time.Now().AddDate(0, 0, -1*maxDays)
	files, err := filepath.Glob(filepath.Join(dir, "news-????-??-??.html"))
	if err != nil {
		return err
	}

	for _, f := range files {
		s := reFile.FindStringSubmatch(f)
		if len(s) != 2 {
			continue
		}
		t, err := time.Parse("2006-01-02", s[1])
		if err != nil {
			continue
		}
		if t.After(cutoff) {
			continue
		}
		if err := os.Remove(f); err != nil {
			return err
		}
		debug("Removed %s", f)
	}
	return nil
}
