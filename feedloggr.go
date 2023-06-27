package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/lmas/feedloggr/internal"
)

type command struct {
	Cmd  string
	Help string
	Func func([]string)
}

var (
	confFile    = flag.String("conf", ".feedloggr.yml", "Path to conf file")
	confVerbose = flag.Bool("verbose", false, "Print debug messages while running")

	commands []command
)

func main() {
	commands = []command{
		{"discover", "Try discover feeds from <URL>", cmdDiscover},
		{"example", "Print example config", cmdExample},
		{"help", "Print this help message and exit", cmdHelp},
		{"regexp", "Try parsing items from <URL> using <regexp> rule", cmdRegexp},
		{"run", "Update feeds and output new page", cmdRun},
		{"test", "Try loading config", cmdTest},
		{"version", "Print version information", cmdVersion},
	}

	flag.Usage = printUsage
	flag.Parse()
	cmd := strings.ToLower(flag.Arg(0))
	args := flag.Args()
	if len(args) > 0 {
		args = args[1:] // Removes the cmd arg
	}

	for _, c := range commands {
		if c.Cmd == cmd {
			c.Func(args)
			return
		}
	}

	printUsage()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func printUsage() {
	out := flag.CommandLine.Output()
	fmt.Fprintf(out, "Usage of %s:\n\n", os.Args[0])
	fmt.Fprintln(out, "Flags")
	flag.PrintDefaults()
	fmt.Fprintln(out, "\nCommands")
	for _, c := range commands {
		fmt.Fprintf(out, "  %s\n\t%s\n", c.Cmd, c.Help)
	}
}

func cmdHelp(args []string) {
	printUsage()
}

func cmdVersion(args []string) {
	// This is supposed to be a toilet/paper roll
	fmt.Printf("  ,-. \n"+
		" ( O )`~-~-~-~-~-~-~-~-~-, \n"+
		" |`-'|  -- %s --\t | \n"+
		" |   |     %s\t | \n"+
		"  `-' `~-~-~-~-~-~-~-~-~-' \n", internal.GeneratorName, internal.GeneratorVersion)
}

func cmdExample(args []string) {
	fmt.Println(internal.ExampleConf())
}

func cmdTest(args []string) {
	conf, err := internal.LoadConf(*confFile)
	if err != nil {
		fmt.Printf("Error loading conf %s: %s\n", *confFile, err)
		return
	}
	fmt.Println(conf)
	fmt.Printf("No errors while loading: %s\n", *confFile)
}

func cmdDiscover(args []string) {
	if len(args) != 1 {
		fmt.Printf("Error discover command expects a single argument: URL, but got: %s\n", args)
		return
	}

	url := strings.ToLower(args[0])
	feeds, err := internal.DiscoverFeeds(url)
	if err != nil {
		fmt.Printf("Error discovering feeds at %s: %s\n", url, err)
		return
	} else if len(feeds) < 1 {
		fmt.Println("No feeds found")
	} else {
		fmt.Println("Possible feeds:")
		for i, f := range feeds {
			fmt.Printf("#%d\t %s\n", i+1, f)
		}
	}
}

func cmdRegexp(args []string) {
	if len(args) != 2 {
		fmt.Printf("Error regexp command expects two arguments: URL, regexp, but got: %s\n", args)
		return
	}
	u, err := url.Parse(args[0])
	if err != nil {
		fmt.Printf("Error parsing url %s: %s\n", args[0], err)
		return
	}

	gen, err := internal.NewGenerator(internal.Conf{
		Settings: internal.Settings{
			Timeout:  10,
			Jitter:   0,
			MaxItems: 30,
		},
	})
	if err != nil {
		fmt.Printf("Error creating generator: %s\n", err)
		return
	}

	items, err := gen.FetchItems(internal.Feed{
		Url: u.String(),
		Parser: internal.Parser{
			Rule: args[1],
			Host: u.Host,
		},
	})
	if err != nil {
		// An error will be returned when the regexp fails to match any items, too
		fmt.Printf("Error fetching items from %s: %s\n", u.String(), err)
		return
	}

	fmt.Println("Items found:")
	for i, item := range items {
		fmt.Printf("#%d\t %s\t (%s)\n", i, item.Title, item.Url)
	}
}

func cmdRun(args []string) {
	conf, err := internal.LoadConf(*confFile)
	if err != nil {
		fmt.Printf("Error loading config %s: %s\n", *confFile, err)
		return
	}

	if *confVerbose != conf.Settings.Verbose {
		// Weeell if one of 'em is true dey bath gotta be true nao
		*confVerbose, conf.Settings.Verbose = true, true
	}
	debug("Loaded config from: %s", *confFile)

	tmpl, err := internal.LoadTemplate(conf.Settings.Template)
	if err != nil {
		fmt.Printf("Error loading template %s: %s\n", conf.Settings.Template, err)
		return
	}

	feeds, err := fetchFeeds(conf)
	if err != nil {
		fmt.Printf("Error fetching feeds: %s\n", err)
		return
	}

	if err := writeFiles(conf.Settings.Output, feeds, tmpl); err != nil {
		fmt.Printf("Error writing files: %s\n", err)
		return
	}

	if err := removeOldFiles(conf.Settings.Output, conf.Settings.MaxDays); err != nil {
		fmt.Printf("Error removing old files: %s\n", err)
		return
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

	for _, feed := range conf.Feeds {
		debug("Updating %s (%s)", feed.Title, feed.Url)
		items, errFeed := gen.NewItems(feed)
		if errFeed != nil {
			debug("\tError: %s", errFeed)
		} else if len(items) > 0 {
			debug("\tItems: %d", len(items))
		} else {
			debug("No items/errors")
			continue
		}

		feeds = append(feeds, internal.TemplateFeed{
			Conf:  feed,
			Items: items,
			Error: errFeed,
		})
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
