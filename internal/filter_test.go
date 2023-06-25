package internal

import (
	"encoding/binary"
	"math/rand"
	"testing"
)

var testItems = []Item{
	{"bbb", "https://bbb.com", ""},
	{"aaa", "https://aaa.com", ""},
	{"ddd", "https://bbb.com", ""},
	{"ccc", "https://ccc.com", ""},
}

func TestSimpleLoadAndFilter(t *testing.T) {
	f, err := loadFilter("")
	if err != nil {
		t.Fatalf("loadfilter: %s", err)
	}
	items := f.filterItems(3, testItems...)
	expected := []Item{
		{"aaa", "https://aaa.com", ""},
		{"bbb", "https://bbb.com", ""},
	}
	if len(items) != len(expected) {
		t.Fatalf("expected %d items, got %d", len(expected), len(items))
	}
	for i := range items {
		if items[i] != expected[i] {
			t.Fatalf("expected item %q, got %q", expected[i], items[i])
		}
	}
	items2 := f.filterItems(3, testItems...)
	if len(items2) != 0 {
		t.Fatalf("expected no items, got %d", len(items2))
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func randomString() string {
	const (
		alphanum string = "abcdefghijklmnopqrstuvwxyz0123456789-_"
		min, max int    = 10, 100
	)
	size := rand.Intn(max-min+1) + min
	b := make([]byte, size)
	for i := range b {
		b[i] = alphanum[rand.Intn(len(alphanum))]
	}
	return string(b)
}

func isUnique(s string, items []Item) bool {
	for _, i := range items {
		if i.Title == s {
			return false
		}
	}
	return true
}

func randomItems(amount int) []Item {
	items := make([]Item, amount)
	for i := range items {
		s := randomString()
		if !isUnique(s, items) {
			continue
		}
		items[i] = Item{s, "https://" + s + ".com", ""}
	}
	return items
}

func TestFilterRandomItems(t *testing.T) {
	f, err := loadFilter("")
	if err != nil {
		t.Fatalf("loadfilter: %s", err)
	}
	const maxItems int = 10000
	rand.Seed(int64(maxItems))
	items := randomItems(maxItems)
	i := f.filterItems(maxItems, items...)
	if len(i) > maxItems {
		t.Fatalf("expected %d items, got %d", maxItems, len(i))
	} else if len(i) < maxItems {
		t.Logf("warning: expected %d items, got %d due to false positives", maxItems, len(i))
	}

	// Replace half of the items with new ones
	for i := range items[:maxItems/2] {
		s := randomString()
		if !isUnique(s, items) {
			continue
		}
		items[i] = Item{s, "https://" + s + ".com", ""}
	}
	i = f.filterItems(maxItems, items...)
	if len(i) > maxItems/2 {
		t.Fatalf("expected %d items, got %d", maxItems/2, len(i))
	} else if len(i) < maxItems/2 {
		t.Logf("warning: expected %d items, got %d due to false positives", maxItems/2, len(i))
	}
	t.Logf("filter stats: %+v\n", f.stats())
}

func TestFilterFalsePositiveRate(t *testing.T) {
	f, err := loadFilter("")
	if err != nil {
		t.Fatalf("loadfilter: %s", err)
	}
	const maxItems int = 50000
	rand.Seed(int64(maxItems))
	buf := make([]byte, 4)
	for i := 0; i < maxItems; i++ {
		binary.BigEndian.PutUint32(buf, uint32(i))
		f.bloom.Add(buf)
	}
	fp := 0
	for i := 0; i < maxItems; i++ {
		binary.BigEndian.PutUint32(buf, uint32(i+maxItems+1))
		if f.bloom.Test(buf) {
			fp++
		}
	}
	ratio := float64(fp) / float64(maxItems)
	if ratio > defaultFilterRate {
		t.Logf("warning: expected ratio <= %f, got %f", defaultFilterRate, ratio)
	} else {
		t.Logf("false positives rate: %f", ratio)
	}
	t.Logf("filter stats: %+v\n", f.stats())
}
