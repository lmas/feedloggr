package main

import (
	"flag"
	"fmt"
	"html/template"
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
	Func func([]string) error
}

var (
	confFile    = flag.String("conf", ".feedloggr.yml", "Path to conf file")
	confVerbose = flag.Bool("verbose", false, "Print debug messages while running")

	commands []command
)

func main() {
	commands = []command{
		{"discover", "Try discover feeds from URL", cmdDiscover},
		{"example", "Print example config", cmdExample},
		{"help", "Print this help message and exit", cmdHelp},
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
			if err := c.Func(args); err != nil {
				panic(err)
			}
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

func cmdHelp(args []string) error {
	printUsage()
	return nil
}

func cmdVersion(args []string) error {
	fmt.Printf("%s %s\n", internal.GeneratorName, internal.GeneratorVersion)
	return nil
}

func cmdExample(args []string) error {
	fmt.Println(internal.ExampleConf())
	return nil
}

func cmdTest(args []string) error {
	conf, err := internal.LoadConf(*confFile)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf)
	fmt.Println("No errors while loading config")
	return nil
}

func cmdDiscover(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("discover command expects a single argument: URL, but got: %s", args)
	}

	url := strings.ToLower(args[0])
	feeds, err := internal.DiscoverFeeds(url)
	if err != nil {
		fmt.Println("Error: ", err)
	} else if len(feeds) < 1 {
		fmt.Println("No feeds found")
	} else {
		fmt.Println("Possible feeds:")
		for i, f := range feeds {
			fmt.Printf("#%d\t %s\n", i+1, f)
		}
	}
	return nil
}

func cmdRun(args []string) error {
	conf, err := internal.LoadConf(*confFile)
	if err != nil {
		return err
	}

	if *confVerbose != conf.Settings.Verbose {
		// Weeell if one of 'em is true dey bath gotta be true nao
		*confVerbose, conf.Settings.Verbose = true, true
	}

	tmpl, err := internal.LoadTemplate(conf.Settings.Template)
	if err != nil {
		return err
	}

	feeds, err := fetchFeeds(conf)
	if err != nil {
		return err
	}

	if err := writeFiles(conf.Settings.Output, feeds, tmpl); err != nil {
		return err
	}

	if err := removeOldFiles(conf.Settings.Output, conf.Settings.MaxDays); err != nil {
		return err
	}
	return nil
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
