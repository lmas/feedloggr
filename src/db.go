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

func (db *DB) SaveItems(items []*FeedItem) {
	tx := db.Begin()
	tx.LogMode(false) // Don't show errors when UNIQUE fails

	for _, i := range items {
		tx.Create(i)
	}

	tx.Commit()
}

func (db *DB) GetItems(feed_url string) []*FeedItem {
	var items []*FeedItem
	// TODO: fix the feed url thing
	db.Order("title, date desc").Where(
		"feed = ? AND date(date) = date(?)", feed_url, Now(),
	).Find(&items)
	return items
}
