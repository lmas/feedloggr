package feedloggr2

import (
	"github.com/jinzhu/gorm"
	rss "github.com/jteeuwen/go-pkg-rss"
	_ "github.com/mattn/go-sqlite3"
)

type Datastore interface {
	GetItems(feed_url string) ([]*FeedItem, error)

	ProcessChannels(feed *rss.Feed, channels []*rss.Channel)
	ProcessItems(feed *rss.Feed, ch *rss.Channel, items []*rss.Item)
}

type DB struct {
	*gorm.DB
}

func OpenSqliteDB(args ...interface{}) (*DB, error) {
	// TODO: get rid of gorm
	db, e := gorm.Open("sqlite3", args...)
	if e != nil {
		return nil, e
	}
	db.AutoMigrate(&FeedItem{})
	return &DB{&db}, nil
}

func (db *DB) GetItems(feed_url string) ([]*FeedItem, error) {
	var items []*FeedItem
	// TODO: fix the feed url thing
	db.Order("date desc, title").Where(
		"feed = ? AND date(date) = date(?)", feed_url, Now(),
	).Find(&items)
	return items, nil
}

// Dummy func so go-pkg-rss will run.
func (db *DB) ProcessChannels(feed *rss.Feed, channels []*rss.Channel) {
}

func (db *DB) ProcessItems(feed *rss.Feed, ch *rss.Channel, items []*rss.Item) {
	tx := db.Begin()
	tx.LogMode(false) // Don't show errors when UNIQUE fails
	for _, it := range items {
		tx.Create(&FeedItem{
			Title: it.Title,
			Url:   it.Links[0].Href,
			Date:  Now(),
			Feed:  feed.Url, // TODO: fix the feed url thing
		})
	}
	tx.Commit()
}
