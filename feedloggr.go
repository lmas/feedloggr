package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lmas/feedloggr/internal"
)

var (
	confFile    = flag.String("conf", ".feedloggr.yaml", "Path to conf file")
	confClean   = flag.Bool("clean", false, "Clean up old pages and exit")
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
	tmpls, err := internal.LoadTemplates()
	if err != nil {
		panic(err)
	}
	switch {
	case *confClean:
		// TODO
		os.Exit(0)
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
	}

	p := filepath.Join(conf.Settings.Output, "news-"+vars.Today.Format("2006-01-02")+".html")
	if err := internal.WriteTemplate(p, "layout.html", tmpls, vars); err != nil {
		panic(err)
	}
	debug("Wrote %s", p)
	if err := internal.Symlink(p, filepath.Join(conf.Settings.Output, "index.html")); err != nil {
		panic(err)
	}

	p = filepath.Join(conf.Settings.Output, "style.css")
	if _, err := os.Stat(p); err != nil {
		// // TODO: should probably make sure the error is an IsNotExist
		if err := internal.WriteTemplate(p, "style.css", tmpls, vars); err != nil {
			panic(err)
		}
		debug("Wrote %s", p)
	}

	if err := gen.WriteFilter(conf.Settings.Output); err != nil {
		panic(err)
	}
	if conf.Settings.Verbose {
		s := gen.FilterStats()
		fmt.Println("Filter Stats")
		fmt.Println("Cells:		", s.Cells)
		fmt.Println("Hashes:		", s.HashFunctions)
		fmt.Println("CellDecrement:	", s.CellDecrement)
		fmt.Println("ActualFalsePos:	", s.FalsePositiveRate)
		fmt.Println("StablePoint:	", s.StablePoint)
	}
}
