package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/lmas/feedloggr2/src"
)

func main() {
	app := cli.NewApp()
	app.Version = feedloggr2.VERSION
	app.Usage = "Collect news from RSS/Atom feeds and create static news web pages."
	app.Flags = []cli.Flag{
		// TODO: default to BoolFlag when done debugging
		cli.BoolTFlag{
			Name:  "verbose",
			Usage: "run in verbose mode",
		},
		cli.StringFlag{
			Name:  "config",
			Usage: "path to config file",
			Value: ".feedloggr2.conf",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "config",
			Usage:  "Print an example config",
			Action: example_config,
		},
		{
			Name:   "test",
			Usage:  "Test the config file",
			Action: test_config,
		},
		{
			Name:   "run",
			Usage:  "Run the generator",
			Action: run,
		},
	}
	app.Run(os.Args)
}

func example_config(c *cli.Context) {
	conf := feedloggr2.NewConfig()
	pretty, e := json.MarshalIndent(conf, "", "    ")
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println(string(pretty))
}

func test_config(c *cli.Context) {
	_, e := feedloggr2.LoadConfig(c.GlobalString("config"))
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	if c.GlobalBool("verbose") {
		fmt.Println("No errors while loading config file.")
	}
}

func run(c *cli.Context) {
	config, e := feedloggr2.LoadConfig(c.GlobalString("config"))
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	e = feedloggr2.UpdateFeeds(config)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}
