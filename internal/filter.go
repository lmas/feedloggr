package internal

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"hash"
	"os"
	"path/filepath"
	"sort"

	boom "github.com/tylertreat/BoomFilters"
)

const (
	// Estimated amount of news items from 50 feeds, with 50 daily updates each, for 365 days = 50*50*365 = 912500
	// defaultFilterSize uint    = 912500 // TODO: items != size of filter
	defaultFilterSize uint    = 1000000
	defaultFilterRate float64 = 0.001
	defaultFilterPath string  = ".filter.dat"
)

type hashwrap struct {
	hash.Hash
}

func (h hashwrap) Sum64() uint64 {
	b := h.Sum(nil)
	return binary.LittleEndian.Uint64(b[:8])
}

func loadFilter(path string) (*boom.StableBloomFilter, error) {
	filter := boom.NewDefaultStableBloomFilter(defaultFilterSize, defaultFilterRate)
	// The default hash is fast but might not be as close to the false positive rate as we expect,
	// so instead use sha256 (slower but more accurate, see https://github.com/tylertreat/BoomFilters/pull/1)
	filter.SetHash(hashwrap{Hash: sha256.New()})
	f, err := os.Open(path) // #nosec G304
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return filter, nil
		}
		return nil, err
	}
	defer f.Close() // #nosec G307
	if _, err := filter.ReadFrom(f); err != nil {
		return nil, err
	}
	return filter, nil
}

func saveFilter(dir string, filter *boom.StableBloomFilter) error {
	path := filepath.Join(dir, defaultFilterPath)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close() // #nosec G307
	_, err = filter.WriteTo(f)
	return err
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// WriteFilter writes the internal Bloom filter to dir
func (g *Generator) WriteFilter(dir string) error {
	return saveFilter(dir, g.filter)
}

// Filter filters out old, already seen items, sorts them by their titles and then cuts of excess amounts
func (g *Generator) Filter(items ...Item) []Item {
	// TODO: avoid making new slice?
	var filtered []Item
	for _, i := range items {
		if !g.filter.TestAndAdd([]byte(i.Url)) {
			filtered = append(filtered, i)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Title < filtered[j].Title
	})
	if len(filtered) > g.conf.Settings.MaxItems {
		filtered = filtered[:g.conf.Settings.MaxItems]
	}
	return filtered
}

// FilterStats contains basic info about the internal Bloom Filter
type FilterStats struct {
	Cells             uint    // number of cells in the stable bloom filter
	CellDecrement     uint    // number of cells decremented on every add
	HashFunctions     uint    // number of hash functions
	FalsePositiveRate float64 // upper bound on false positives when the filter has become stable
	// the limit of the expected fraction of zeros in the filter when the number of iterations
	// goes to infinity. When this limit is reached, the filter is considered stable.
	StablePoint float64
}

// FilterStats returns basic info about the internal Bloom Filter...
func (g *Generator) FilterStats() FilterStats {
	return FilterStats{
		Cells:             g.filter.Cells(),
		FalsePositiveRate: g.filter.FalsePositiveRate(),
		HashFunctions:     g.filter.K(),
		CellDecrement:     g.filter.P(),
		StablePoint:       g.filter.StablePoint(),
	}
}
