package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lmas/feedloggr2"
)

var (
	verbose = flag.Bool("verbose", false, "run in verbose mode")
	config  = flag.String("config", ".feedloggr2.conf", "path to config file")

	version = flag.Bool("version", false, "print version and exit")
	example = flag.Bool("example", false, "print example config and exit")
	test    = flag.Bool("test", false, "test config file and exit")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("feedloggr2 v%s\n", feedloggr2.Version)
		fmt.Println("Collect news from RSS/Atom feeds and create static news pages in HTML.")
		return
	}

	if *example {
		c := feedloggr2.NewConfig()
		fmt.Println(c)
		return // simple exit(0)
	}

	c, err := feedloggr2.LoadConfig(*config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *test {
		fmt.Println(c)
		fmt.Println("No errors while loading config file.")
		return
	}

	// cmd flags override config file.
	if *verbose {
		c.Verbose = true
	}

	err = feedloggr2.UpdateFeeds(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
