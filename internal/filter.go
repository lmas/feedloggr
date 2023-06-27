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

// FilterStats contains basic info about the internal Bloom Filter.
type FilterStats struct {
	Capacity  uint    // Total capacity for the internal series of bloom filters
	Hashes    uint    // Number of hash functions for each internal filter
	FillRatio float64 // Average ratio of set bits across all internal filters
}

// String returns a pretty string of FilterStats.
func (fs FilterStats) String() string {
	return fmt.Sprintf("Capacity = %d, Hashes = %d, Fill Ratio = %f",
		fs.Capacity,
		fs.Hashes,
		fs.FillRatio,
	)
}

// returns basic info for a Bloom Filter.
func (f *filter) stats() FilterStats {
	return FilterStats{
		f.bloom.Capacity(),
		f.bloom.K(),
		f.bloom.FillRatio(),
	}
}

////////////////////////////////////////////////////////////////////////////////

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
	bloom *boom.ScalableBloomFilter
	path  string
}

func loadFilter(dir string) (*filter, error) {
	bloom := boom.NewDefaultScalableBloomFilter(defaultFilterRate)
	// The default hash is fast but might not be as close to the false positive rate as we expect,
	// so instead use sha256 (slower but more accurate, see https://github.com/tylertreat/BoomFilters/pull/1)
	bloom.SetHash(hashwrap{Hash: sha256.New()})
	filter := &filter{
		bloom: bloom,
		path:  filepath.Join(dir, defaultFilterPath),
	}

	f, err := os.Open(filter.path) // #nosec G304
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// If not done for a new and empty filter, the filter data won't be
			// saved properly after the first run for some reason
			// (might be upstream bug)
			bloom = bloom.Reset()
			return filter, nil
		}
		return nil, err
	}

	defer f.Close() // #nosec G307
	if _, err := bloom.ReadFrom(f); err != nil {
		return nil, err
	}
	return filter, nil
}

func (f *filter) write() error {
	fd, err := os.Create(f.path)
	if err != nil {
		return err
	}

	defer fd.Close() // #nosec G307
	_, err = f.bloom.WriteTo(fd)
	return err
}

func (f *filter) filterItems(max int, items ...Item) []Item {
	if len(items) < 1 {
		return nil
	}

	var filtered []Item
	for _, i := range items {
		if f.bloom.TestAndAdd([]byte(i.Url)) == false {
			filtered = append(filtered, i)
			if len(filtered) >= max {
				break
			}
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Title < filtered[j].Title
	})
	return filtered
}
