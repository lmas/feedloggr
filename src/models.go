package feedloggr2

import "time"

type Feed struct {
	Title string
	Url   string
	Items []*FeedItem
}

// TODO: make indexes
type FeedItem struct {
	ID    int
	Title string `sql:"type:varchar(100)"`
	Url   string `sql:"type:varchar(255);unique"`
	Date  time.Time
	Feed  string `sql:"type:varchar(64)"`
}
