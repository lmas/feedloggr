package feedloggr2

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const sqlCreateTable string = `CREATE TABLE IF NOT EXISTS feed_items (
	id INTEGER PRIMARY KEY,
	title TEXT,
	url TEXT UNIQUE,
	date DATE,
	feed TEXT
);
CREATE INDEX IF NOT EXISTS index_feed_items ON feed_items(feed, date);
`

const sqlInsertItem string = `INSERT OR IGNORE INTO feed_items VALUES(
	NULL,
	?,
	?,
	?,
	?
);
`

const sqlGetItems string = `SELECT title, url FROM feed_items
	WHERE feed = ? AND date = ?
	ORDER BY title, date DESC;
`

type DB struct {
	*sqlx.DB
}

func OpenDB(path string) (*DB, error) {
	conn, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(sqlCreateTable)
	if err != nil {
		return nil, err
	}

	db := &DB{conn}
	return db, nil
}

func (db *DB) SaveItems(feedURL, date string, items []Item) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, i := range items {
		if i.Title == "" || i.URL == "" {
			// Yeah don't want empty/weird stuff getting stuck in the db
			continue
		}

		_, err := tx.Exec(sqlInsertItem, i.Title, i.URL, date, feedURL)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetItems(feedURL, date string) ([]Item, error) {
	var items []Item
	err := db.Select(&items, sqlGetItems, feedURL, date)
	if err != nil {
		return nil, err
	}
	return items, nil
}
