/*
Copyright © 2021 James Lucktaylor <jlucktay@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

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
func benchmarkSortByAirdate(b *testing.B, limit int, sortFn func(data sort.Interface)) {
	b.Helper()

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

func BenchmarkSortByAirdateSort32(b *testing.B)  { benchmarkSortByAirdate(b, 32, sort.Sort) }
func BenchmarkSortByAirdateSort64(b *testing.B)  { benchmarkSortByAirdate(b, 64, sort.Sort) }
func BenchmarkSortByAirdateSort128(b *testing.B) { benchmarkSortByAirdate(b, 128, sort.Sort) }
func BenchmarkSortByAirdateSort256(b *testing.B) { benchmarkSortByAirdate(b, 256, sort.Sort) }
func BenchmarkSortByAirdateSort512(b *testing.B) { benchmarkSortByAirdate(b, 512, sort.Sort) }

func BenchmarkSortByAirdateStable32(b *testing.B)  { benchmarkSortByAirdate(b, 32, sort.Stable) }
func BenchmarkSortByAirdateStable64(b *testing.B)  { benchmarkSortByAirdate(b, 64, sort.Stable) }
func BenchmarkSortByAirdateStable128(b *testing.B) { benchmarkSortByAirdate(b, 128, sort.Stable) }
func BenchmarkSortByAirdateStable256(b *testing.B) { benchmarkSortByAirdate(b, 256, sort.Stable) }
func BenchmarkSortByAirdateStable512(b *testing.B) { benchmarkSortByAirdate(b, 512, sort.Stable) }
