package inmemory_test

import (
	"fmt"
	"sort"
	"sync"
	"testing"

	"github.com/matryer/is"

	"go.jlucktay.dev/arrowverse/pkg/collection/inmemory"
	"go.jlucktay.dev/arrowverse/pkg/models"
	"go.jlucktay.dev/arrowverse/pkg/scrape"
)

var (
	// result is stored in a package level variable so that the compiler cannot eliminate the Benchmark itself.
	// cf. https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
	result int

	episodes   []models.Episode
	scrapeOnce sync.Once
)

func populateCollection() error {
	var err error

	scrapeOnce.Do(func() {
		var episodeLists map[models.ShowName]string
		episodes = make([]models.Episode, 0)

		if episodeLists, err = scrape.EpisodeLists(); err != nil {
			return
		}

		for s, el := range episodeLists {
			var show *models.Show

			if show, err = scrape.Episodes(s, el); err != nil {
				return
			}

			for seasonIndex := range show.Seasons {
				episodes = append(episodes, show.Seasons[seasonIndex].Episodes...)
			}
		}
	})

	if err != nil {
		return fmt.Errorf("could not populate collection: %w", err)
	}

	return nil
}

// benchmarkSortByAirdate is private to only be invoked by wrappers requesting different size data sets.
func benchmarkSortByAirdate(limit int, sortFn func(data sort.Interface), b *testing.B) {
	var r int

	is := is.New(b)

	is.NoErr(populateCollection())
	is.True(len(episodes) >= limit)

	limitedSet := episodes[:limit]

	b.ResetTimer() // If a benchmark needs some expensive setup before running, the timer may be reset

	for n := 0; n < b.N; n++ {
		sortFn(inmemory.ByAirdate(limitedSet))

		r = len(limitedSet[0].String())
	}

	result = r
}

func BenchmarkSortByAirdateSort32(b *testing.B)  { benchmarkSortByAirdate(32, sort.Sort, b) }
func BenchmarkSortByAirdateSort64(b *testing.B)  { benchmarkSortByAirdate(64, sort.Sort, b) }
func BenchmarkSortByAirdateSort128(b *testing.B) { benchmarkSortByAirdate(128, sort.Sort, b) }
func BenchmarkSortByAirdateSort256(b *testing.B) { benchmarkSortByAirdate(256, sort.Sort, b) }
func BenchmarkSortByAirdateSort512(b *testing.B) { benchmarkSortByAirdate(512, sort.Sort, b) }

func BenchmarkSortByAirdateStable32(b *testing.B)  { benchmarkSortByAirdate(32, sort.Stable, b) }
func BenchmarkSortByAirdateStable64(b *testing.B)  { benchmarkSortByAirdate(64, sort.Stable, b) }
func BenchmarkSortByAirdateStable128(b *testing.B) { benchmarkSortByAirdate(128, sort.Stable, b) }
func BenchmarkSortByAirdateStable256(b *testing.B) { benchmarkSortByAirdate(256, sort.Stable, b) }
func BenchmarkSortByAirdateStable512(b *testing.B) { benchmarkSortByAirdate(512, sort.Stable, b) }
