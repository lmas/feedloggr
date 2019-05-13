package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lmas/feedloggr"
)

var (
	verbose = flag.Bool("verbose", false, "run in verbose mode")
	config  = flag.String("config", ".feedloggr3.conf", "path to config file")

	version = flag.Bool("version", false, "print version and exit")
	example = flag.Bool("example", false, "print example config and exit")
	test    = flag.Bool("test", false, "test config file and exit")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Println("Feedloggr 3.0")
		fmt.Println("Aggregate news from RSS/Atom feeds and output static HTML pages.")
		return
	}

	if *example {
		cfg := feedloggr.NewConfig()
		fmt.Println(cfg)
		return // simple exit(0)
	}

	cfg, err := feedloggr.LoadConfig(*config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *test {
		fmt.Println(cfg)
		fmt.Println("No errors while loading config file.")
		return
	}

	// cmd flags override config file
	if *verbose {
		cfg.Verbose = true
	}

	app, err := feedloggr.New(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = app.Update()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
