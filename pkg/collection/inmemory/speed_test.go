package inmemory_test

import (
	"sync"
	"testing"

	"github.com/matryer/is"

	"go.jlucktay.dev/arrowverse/pkg/collection"
	"go.jlucktay.dev/arrowverse/pkg/collection/inmemory"
	"go.jlucktay.dev/arrowverse/pkg/models"
	"go.jlucktay.dev/arrowverse/pkg/scrape"
)

var (
	// result is stored in a package level variable so that the compiler cannot eliminate the Benchmark itself.
	// cf. https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
	result int

	cs         collection.Shows
	scrapeOnce sync.Once
)

func getEpisodes() (collection.Shows, error) {
	var err error

	scrapeOnce.Do(func() {
		var episodeLists map[models.ShowName]string
		cs = &inmemory.Collection{}

		if episodeLists, err = scrape.EpisodeLists(); err != nil {
			return
		}

		for s, el := range episodeLists {
			var show *models.Show

			if show, err = scrape.Episodes(s, el); err != nil {
				return
			}

			if err = cs.Add(show); err != nil {
				return
			}
		}
	})

	if err != nil {
		return nil, err
	}

	return cs, nil
}

// benchmarkInOrder is private so as not to be invoked directly, but by wrappers requesting different size data sets.
// TODO:
// - wire up the int parameter to actually do the different size data sets thing
// - bench sort.Sort against sort.Stable inside the .InOrder() call
func benchmarkInOrder(_ int, b *testing.B) {
	is := is.New(b)

	episodes, errGE := getEpisodes()
	is.NoErr(errGE)

	var r int

	b.ResetTimer() // If a benchmark needs some expensive setup before running, the timer may be reset

	for n := 0; n < b.N; n++ {
		episodesInOrder, errIO := episodes.InOrder()
		is.NoErr(errIO)

		r = len(episodesInOrder)
	}

	result = r
}

func BenchmarkInOrder1(b *testing.B)  { benchmarkInOrder(1, b) }
func BenchmarkInOrder2(b *testing.B)  { benchmarkInOrder(2, b) }
func BenchmarkInOrder4(b *testing.B)  { benchmarkInOrder(4, b) }
func BenchmarkInOrder8(b *testing.B)  { benchmarkInOrder(8, b) }
func BenchmarkInOrder16(b *testing.B) { benchmarkInOrder(16, b) }
