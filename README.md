
# feedloggr v0.3
[![PkgGoDev](https://pkg.go.dev/badge/github.com/lmas/feedloggr)](https://pkg.go.dev/github.com/lmas/feedloggr)

Collect news from your favourite Atom/RSS/JSON feeds and generate static web pages for easy browsing.


## Status

The project has successfully been running since June 2021 without any major issues on my own sites.
And the only minor issues discovered so far has been with external feeds with poor behaviour.


## Installation

    go install github.com/lmas/feedloggr


## Usage

For a list of available flags and commands:

    feedloggr -h

You can create a new configuration file by running:

    feedloggr -example > .feedloggr.yml

You should then edit file and add in your news feeds.

When you're done editing the configuration, you can test it and make sure there's no errors in it:

    feedloggr -test

If no errors are shown, you're good to go.

Now that you have a working configuration,
it's time to start collect news from your feeds and create a new web page showing all the collected news:

    feedloggr

When it's done you should be able to browse the newly generated pages,
found inside the output directory that was specified in the configuration file.


## Command line options

```
Usage of feedloggr:
    -clean
        Clean up old pages and exit
    -conf string
        Path to conf file (default ".feedloggr.yml")
    -example
        Print example config and exit
    -test
        Load config and exit
    -verbose
        Print debug messages while running
    -version
        Print version and exit
```


## Configuration

Configuration is by default loaded from the file `.feedloggr.yml` in the current directory,
but can be overridden with the `-conf` flag.

Example configuration file:

    settings:
      output: ./feeds/
      template: new.html
      maxitems: 20
      throttle: 2
      timeout: 30
      verbose: true
    feeds:
    - title: Reddit
      url: https://old.reddit.com/.rss
    - title: Tech - Hacker News
      url: https://news.ycombinator.com/rss
      parser:
        rule: (?sU)<item>.*<title>(?P<title>[^<]+)</title>.*<comments>(?P<url>[^<]+)</comments>.*</item>
        host: https://news.ycombinator.com/rss

### settings

Global configuration settings.

**output**

    Output directory where generated pages and link filter are stored.

**template**

    Optional filepath to custom HTML template.

**maxitems**

    Max amount of items to fetch from each feed.

**throttle**

    Stupid simple time to sleep (in seconds) between each feed update.

**timeout**

    Max time (in seconds) to try download feed before timing out.

**verbose**

   Show verbose output when running.

### feeds

List of Atom/RSS/JSON feeds.

**title**

    Custom title for a feed.

**url**

    Source feed URL for fetching new updates.

### parser

Regexp rules to either parse a non-standard feed or for fetching specific parts of a regular feed.

**rule**

    Regexp rule for fetching items from source URL. It must provide two capture groups called "title" and "url".

**host**

    Optional host value for item URLs with a missing host prefix.


## License

GPL, See the [LICENSE](LICENSE) file for details.

