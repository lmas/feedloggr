package internal

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"sort"

	boom "github.com/tylertreat/BoomFilters"
)

const (
	defaultFilterRate float64 = 0.0001
	defaultFilterPath string  = ".filter.dat"
)

type hashwrap struct {
	hash.Hash
}

func (h hashwrap) Sum64() uint64 {
	b := h.Sum(nil)
	return binary.LittleEndian.Uint64(b[:8])
}

type filter struct {
	b *boom.ScalableBloomFilter
}

func loadFilter(path string) (*filter, error) {
	b := boom.NewDefaultScalableBloomFilter(defaultFilterRate)
	// The default hash is fast but might not be as close to the false positive rate as we expect,
	// so instead use sha256 (slower but more accurate, see https://github.com/tylertreat/BoomFilters/pull/1)
	b.SetHash(hashwrap{Hash: sha256.New()})
	filter := &filter{
		b: b,
	}
	f, err := os.Open(path) // #nosec G304
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return filter, nil
		}
		return nil, err
	}
	defer f.Close() // #nosec G307
	if _, err := filter.b.ReadFrom(f); err != nil {
		return nil, err
	}
	return filter, nil
}

func (f *filter) write(dir string) error {
	path := filepath.Join(dir, defaultFilterPath)
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close() // #nosec G307
	_, err = f.b.WriteTo(fd)
	return err
}

func (f *filter) filterItems(max int, items ...Item) []Item {
	if len(items) < 1 {
		return nil
	} else if len(items) > max {
		items = items[:max]
	}
	// TODO: avoid making new slice?
	var filtered []Item
	for _, i := range items {
		if f.b.TestAndAdd([]byte(i.Url)) == false {
			filtered = append(filtered, i)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Title < filtered[j].Title
	})
	return filtered
}

// FilterStats contains basic info about the internal Bloom Filter
type FilterStats struct {
	Capacity  uint    // Total capacity for the internal series of bloom filters
	Hashes    uint    // Number of hash functions for each internal filter
	FillRatio float64 // Average ratio of set bits across all internal filters
}

// String returns a pretty string of FilterStats
func (fs FilterStats) String() string {
	return fmt.Sprintf("Capacity = %d, Hashes = %d, Fill Ratio = %f",
		fs.Capacity,
		fs.Hashes,
		fs.FillRatio,
	)
}

// FilterStats returns basic info about the internal Bloom Filter...
func (f *filter) stats() FilterStats {
	return FilterStats{
		f.b.Capacity(),
		f.b.K(),
		f.b.FillRatio(),
	}
}
