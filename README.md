
     ,-.
    ( O )`~-~-~-~-~-~-~-~-~-,
    |`-'|  -- feedloggr --	 |
    |   |     v0.4.1	 |
     `-' `~-~-~-~-~-~-~-~-~-'

Collect news from your favourite Atom/RSS/JSON feeds and generate static web pages for easy browsing.


## Status

The project has successfully been running since June 2021 without any major issues on my own sites.
The only minor issues discovered so far has been with external feeds with poor behaviour.


## Installation

    go install github.com/lmas/feedloggr@latest


## Usage

For a list of available flags and commands:

    feedloggr help

You can create a new configuration file by running:

    feedloggr example > .feedloggr.yml

You should then edit `.feedloggr.yml` and add your own feeds to it.

When you're done editing the configuration, you can test it and make sure there's no errors:

    feedloggr test

If no errors are shown, you're good to go!

Now that you have a working configuration,
it's time to start collect news from your feeds and export them to a web page:

    feedloggr

When it's done you should be able to browse the newly generated page,
found inside the output directory that was specified in the configuration file.


## Command line options

    Usage of feedloggr:

    Flags
      -conf string
            Path to conf file (default ".feedloggr.yml")
      -verbose
            Print debug messages while running

    Commands
      discover
        Try discover feeds from <URL>
      example
        Print example config
      help
        Print this help message and exit
      regexp
        Try parsing items from <URL> using <regexp> rule
      run
        Update feeds and output new page
      test
        Try loading config
      version
        Print version information


## Configuration

Configuration is by default loaded from the file `.feedloggr.yml` in the current directory,
but can be overridden with the `-conf` flag.

### settings

Global configuration settings.

*output*

    Output directory where generated pages and link filter are stored.

*template*

    Optional filepath to custom HTML template.

*maxdays*

    Optional max amount of days to keep the generated pages for.

*maxitems*

    Max amount of new items to fetch from each feed.

*timeout*

    Max time (in seconds) to try download a feed, before timing out.

*jitter*

    Max time (in seconds) to randomly apply, as a wait time, between each feed update.

*verbose*

    Show verbose output when running.

### feeds

List of Atom/RSS/JSON feeds.

*title*

    Custom title for a feed.

*url*

    Source feed URL for fetching new updates.

*parser.rule*

    Regexp rule for fetching items from a non-feed URL.
    It must provide two capture groups called "title" and "url".
    A third, optional capture group "content" allows for capturing any exra text
    that can be used and displayed in the output template.

*parser.host*

    Optional host prefix for feed item URLs, which can be used to replace a missing
    value or redirect the URLs to another host.

### Example

    settings:
        output: ./feeds/
        template: ""
        maxdays: 30
        maxitems: 20
        timeout: 30
        jitter: 2
        verbose: true
    feeds:
        - title: Lemmy.link
          url: https://lemmy.link/feeds/local.xml?sort=TopDay
        - title: Hacker News
          url: https://news.ycombinator.com/rss
          parser:
            rule: (?sU)<item>.*<title>(?P<title>[^<]+)</title>.*<comments>(?P<url>[^<]+)</comments>.*</item>
            host: https://news.ycombinator.com/rss


## Output Template

The `template` config value can be used to load a custom HTML template,
but if not set it will default to a built in template.

The templating system used is [html/template], found in Go's standard library.

[html/template]: https://pkg.go.dev/html/template

### Template Variables

*.Today*

    Current time. Can be customised with the template functions available down below.

*.Generator.Name*

*.Generator.Version*

*.Generator.Source*

    Basic information about this tool.

*.Feeds*

    A list of feeds, as defined in the config file.
    Can be iterated easily.

*.Feeds.Conf.Title*

*.Feeds.Conf.Url*

    Basic information about current feed.

*.Feeds.Conf.Source*

    In case the `parser.host` setting has been used in the config,
    this variable will simple be the same value.
    Otherwise it defaults to `.Feeds.Conf.Url`.

*.Feeds.Items*

    A list of unique items for the feed.
    Can be iterated easily.

*.Feeds.Items.Title*

*.Feeds.Items.Url*

*.Feeds.Items.Content*

    Basic information about current item.

*.Feeds.Error*

    In case an error was encountered while trying to update the feed,
    this variable will contain the error message.

### Template Functions

*shortdate*

    Can be used to shorten a long time value down to, for example: 2006-01-02.

*prevday*

    Can be used to subtract a day from a time value.

*nextday*

    Can be used to add a day to a time value.

### Example

Minimal example template, based on the built in default (without CSS styling):

    <!doctype html>
    <html>
    <head>
        <meta charset="utf-8">
        <title>News for {{.Today | shortdate}}</title>
        <!-- Optional stylesheet to overwrite the default style -->
        <link href="style.css" rel="stylesheet" type="text/css">
    </head>
    <body>
        <header>
            <h1>{{.Today | shortdate}}</h1>
            <a class="bar" href="./news-{{.Today | prevday | shortdate}}.html">Previous</a>
            <a class="bar" href="./index.html">Latest</a>
            <a class="bar" href="./news-{{.Today | nextday | shortdate}}.html">Next</a>
        </header>
        {{range .Feeds}}
        <section>
            <h2 class="bar"><a href="{{.Conf.Source}}">{{.Conf.Title}}</a></h2>
            {{if .Error}}
                <p>Error while updating feed:<br />{{.Error}}</p>
            {{else}}
                <ul>{{range .Items}}
                    <li><a href="{{.Url}}" rel="nofollow">{{.Title}}</a></li>
                {{end}}</ul>
            {{end}}
        </section>
        {{else}}
            <p>Sorry, no news for today!</p>
        {{end}}
        <footer>
            Generated with <a href="{{.Generator.Source}}">{{.Generator.Name}} {{.Generator.Version}}</a>
        </footer>
    </body>
    </html>

## License

GPL, See the [LICENSE](LICENSE) file for details.

