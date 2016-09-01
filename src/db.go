package feedloggr2

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*gorm.DB
}

func OpenSqliteDB(args ...interface{}) (*DB, error) {
	db, e := gorm.Open("sqlite3", args...)
	if e != nil {
		return nil, e
	}
	db.AutoMigrate(&FeedItem{})
	db.Model(&FeedItem{}).AddIndex("idx_feed_item", "feed", "date")
	return &DB{db}, nil
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
	db.Order("title, date desc").Where(
		"feed = ? AND date(date) = date(?)", feed_url, Now(),
	).Find(&items)
	return items
}
