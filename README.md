
feedloggr2
================================================================================

Collect news from your favorite RSS/Atom feeds and generate simple and static
web pages to browse the news in.

Status
--------------------------------------------------------------------------------

The project is still in an alpha stage and under development.

Installation
--------------------------------------------------------------------------------

    go install github.com/lmas/feedloggr2

Usage
--------------------------------------------------------------------------------

Run `feedloggr2 help` to get a list of available flags and commands.

You can create a new config file by running `feedloggr2 config > .feedloggr2.conf`.
You should then edit `.feedloggr2.conf` and add in your RSS/Atom feeds.
The format of the config is JSON.

When you're done editting the config, you can test it and make sure there's no
errors in it by running `feedloggr2 test`. If no errors are shown, you're good
to go.

Now that you have a working config, it's time to start collect news from your
feeds and create a new web page showing all the collected news.
Run `feedloggr2 run` and it will start working.

When it's done you should be able to browse the newly generated pages, found
inside the output directory.

Configuration
--------------------------------------------------------------------------------

TODO: update when the config values are stable enough.

License
--------------------------------------------------------------------------------

MIT License, see LICENSE.

TODO
--------------------------------------------------------------------------------

Tests:
- Must have unit tests. Need to mock the feed downloading.

Update:
- Set a max on the amount of items gotten.

Config:
- Handle duplicate feeds.
- Make sure the config contains proper values when running the test action.
