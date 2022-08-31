package main

import (
	"flag"
	"fmt"
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

type generator struct {
	Name    string
	Version string
	Source  string
}

type feed struct {
	Conf  internal.Feed
	Items []internal.Item
	Error error
}

type vars struct {
	Generator generator
	Today     time.Time
	Feeds     []feed
}

func debug(msg string, args ...interface{}) {
	if *confVerbose {
		fmt.Printf(msg+"\n", args...)
	}
}

func main() {
	flag.Parse()
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
	switch {
	case *confTest:
		fmt.Println(conf)
		fmt.Println("No errors while loading config")
		os.Exit(0)
	case *confVerbose != conf.Settings.Verbose:
		// Weeell if one of 'em is true dey bath gotta be true nao
		*confVerbose, conf.Settings.Verbose = true, true
	}

	tmpl, err := internal.LoadTemplate(conf.Settings.Template)
	if err != nil {
		panic(err)
	}
	gen, err := internal.New(conf)
	if err != nil {
		panic(err)
	}
	vars := vars{
		Generator: generator{
			Name:    internal.GeneratorName,
			Version: internal.GeneratorVersion,
			Source:  internal.GeneratorSource,
		},
		Today: time.Now(),
	}
	throttle := time.Duration(conf.Settings.Throttle) * time.Second
	for _, source := range conf.Feeds {
		debug("Updating %s (%s)", source.Title, source.Url)
		f := feed{
			Conf: source,
		}
		f.Items, f.Error = gen.NewItems(source)
		if f.Error != nil {
			debug("\tError: %s", f.Error)
		} else {
			debug("\tItems: %d", len(f.Items))
		}
		if len(f.Items) > 0 || f.Error != nil {
			vars.Feeds = append(vars.Feeds, f)
		}
		time.Sleep(throttle)
	}

	p := filepath.Join(conf.Settings.Output, "news-"+vars.Today.Format("2006-01-02")+".html")
	if err := internal.WriteTemplate(p, tmpl, vars); err != nil {
		panic(err)
	}
	debug("Wrote %s", p)
	if err := internal.Symlink(p, filepath.Join(conf.Settings.Output, "index.html")); err != nil {
		panic(err)
	}

	if err := gen.WriteFilter(conf.Settings.Output); err != nil {
		panic(err)
	}
	if conf.Settings.Verbose {
		fmt.Printf("Filter stats: %+v\n", gen.FilterStats())
	}

	if err := removeOldFiles(conf.Settings.Output, conf.Settings.MaxDays); err != nil {
		panic(err)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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
