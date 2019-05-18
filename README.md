
feedloggr v3.0
================================================================================

Collect news from your favorite RSS/Atom feeds and generate simple and static
web pages to browse the news in.

Status
--------------------------------------------------------------------------------

The project is now in beta stage and being actively used at https://lmas.se/news

Installation
--------------------------------------------------------------------------------

    go install github.com/lmas/feedloggr

Usage
--------------------------------------------------------------------------------

For a list of available flags and commands:

    feedloggr -h

You can create a new config file by running:

    feedloggr -example > .feedloggr3.conf

You should then edit `.feedloggr3.conf` and add in your RSS/Atom feeds.
The format of the config is [TOML](https://github.com/toml-lang/toml).

When you're done editting the config, you can test it and make sure there's no
errors in it:

    feedloggr -test

If no errors are shown, you're good to go.

Now that you have a working config, it's time to start collect news from your
feeds and create a new web page showing all the collected news:

    feedloggr

When it's done you should be able to browse the newly generated pages, found
inside the output directory that was specified in the config file.

Configuration
--------------------------------------------------------------------------------

    Verbose = true | false

Show verbose logs while running. Default to false.

    OutputPath = "/path/to/dir"

Path to directory where new pages will be stored in.

    Timeout = seconds

Read timeout when sending HTTP GET requests, when updating the feeds.

    [Feeds]
    "Example" = "https://example.com/rss"
    "Example2" = "https://somewhereelse.com/rss"

List of RSS/Atom feeds, `feedloggr` tries to download the feeds from the URLs.

License
--------------------------------------------------------------------------------

MIT License, see LICENSE.

