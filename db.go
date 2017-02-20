package feedloggr2

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const sqlCreateTable string = `CREATE TABLE IF NOT EXISTS feed_items (
	id INTEGER PRIMARY KEY,
	title TEXT,
	url TEXT NOT NULL UNIQUE,
	date DATE,
	feed TEXT
);

CREATE INDEX IF NOT EXISTS index_feed_item ON feed_items(feed, date);
`

const sqlInsertItem string = `INSERT OR IGNORE INTO feed_items VALUES(
	NULL,
	?,
	?,
	?,
	?
);
`

const sqlGetItems string = `SELECT * FROM feed_items
	WHERE feed = ? AND date(date) = date(?)
	ORDER BY title, date DESC;
`

type DB struct {
	*sqlx.DB
}

func OpenSqliteDB(path string) (*DB, error) {
	conn, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(sqlCreateTable)
	if err != nil {
		return nil, err
	}

	return &DB{conn}, nil
}

func (db *DB) SaveItems(items []*FeedItem) {
	tx, err := db.Begin()
	// TODO: Handle this error better
	if err != nil {
		panic(err)
	}

	now := Now()

	for _, i := range items {
		_, err := tx.Exec(sqlInsertItem, i.Title, i.URL, now, i.Feed)
		// TODO: handle this better
		if err != nil {
			panic(err)
		}
	}

	err = tx.Commit()
	// TODO: Handle this error better
	if err != nil {
		panic(err)
	}
}

func (db *DB) GetItems(feedUrl string) []*FeedItem {
	var items []*FeedItem
	err := db.Select(&items, sqlGetItems, feedUrl, Now())
	// TODO: Handle this error better
	if err != nil {
		panic(err)
	}
	return items
}
