package feedloggr2

import (
	"reflect"
	"testing"
)

func TestDB_OpenSaveItemsGetItems(t *testing.T) {
	items := []Item{
		{"item1", "item1url"},
		{"item1duplicate", "item1url"},
		{"item2", ""},
		{"", "item3url"},
		{"", ""},
		{"item5", "item5url"},
	}
	want := []Item{
		{"item1", "item1url"},
		{"item5", "item5url"},
	}

	db, err := OpenDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to open a sqlite3 database in memory: %s", err)
	}

	err = db.SaveItems("feedurl", "1970-01-01", items)
	if err != nil {
		t.Fatalf("Failed to save items to db: %s", err)
	}

	got, err := db.GetItems("feedurl", "1970-01-01")
	if err != nil {
		t.Fatalf("Failed to get items from db: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("db.GetItems(...) = %v, want %v", got, want)
	}
}
