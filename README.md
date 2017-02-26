
feedloggr2
================================================================================

Collect news from your favorite RSS/Atom feeds and generate simple and static
web pages to browse the news in.

Status
--------------------------------------------------------------------------------

The project is now in beta stage and being actively used at https://lmas.se/news

Installation
--------------------------------------------------------------------------------

    go install github.com/lmas/feedloggr2

Usage
--------------------------------------------------------------------------------

For a list of available flags and commands:

    feedloggr2 -h

You can create a new config file by running:

    feedloggr2 -example > .feedloggr2.conf

You should then edit `.feedloggr2.conf` and add in your RSS/Atom feeds.
The format of the config is [TOML](https://github.com/toml-lang/toml).

When you're done editting the config, you can test it and make sure there's no
errors in it:

    feedloggr2 -test

If no errors are shown, you're good to go.

Now that you have a working config, it's time to start collect news from your
feeds and create a new web page showing all the collected news:

    feedloggr2

When it's done you should be able to browse the newly generated pages, found
inside the output directory that was specified in the config file.

Configuration
--------------------------------------------------------------------------------

    Verbose = true | false

Show verbose logs while running. Default to false.

    Database = "/path/to/file"

Path to file where the sqlite3 database will be written.

    OutputPath = "/path/to/dir"

Path to directory where new pages will be stored in.

    [[Feeds]]
    Title = "Example"
    URL = "https://example.com/rss"

    [[Feeds]]
    Title = "Example2"
    URL = "https://somewhereelse.com/rss"

List of RSS/Atom feeds, `feedloggr2` tries to download the feeds from the URLs.

License
--------------------------------------------------------------------------------

MIT License, see LICENSE.

