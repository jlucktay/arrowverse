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

package scrape_test

import (
	"testing"
	"time"

	"github.com/matryer/is"

	"go.jlucktay.dev/arrowverse/pkg/models"
	"go.jlucktay.dev/arrowverse/pkg/scrape"
)

func TestEpisodeNumbers(t *testing.T) {
	if testing.Short() {
		t.Skip("sometimes takes a little while to scrape")
	}

	t.Parallel()

	is := is.New(t)

	// Store show/season/episode details for seasons that have finished airing completely.
	showSeasonEpisodes := map[models.ShowName]map[int]int{
		models.Arrow:                 {1: 23, 2: 23, 3: 23, 4: 23, 5: 23, 6: 23, 7: 22, 8: 10},
		models.Batwoman:              {1: 20, 2: 18, 3: 13},
		models.BirdsOfPrey:           {1: 13},
		models.BlackLightning:        {1: 13, 2: 16, 3: 16, 4: 13},
		models.Constantine:           {1: 13},
		models.DCsLegendsOfTomorrow:  {1: 16, 2: 17, 3: 18, 4: 16, 5: 15, 6: 15, 7: 13},
		models.FreedomFightersTheRay: {1: 6, 2: 6},
		models.Supergirl:             {1: 20, 2: 22, 3: 23, 4: 22, 5: 19, 6: 20},
		models.TheFlashCBS:           {1: 22},
		models.TheFlashTheCW:         {1: 23, 2: 23, 3: 23, 4: 23, 5: 22, 6: 19, 7: 18},
		models.Vixen:                 {1: 6, 2: 6},
	}

	episodeLists, errEL := scrape.EpisodeLists()
	is.NoErr(errEL)
	is.Equal(len(episodeLists), len(showSeasonEpisodes))

	for s, el := range episodeLists {
		// Pin! ref: https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable
		s, el := s, el

		t.Run(string(s), func(t *testing.T) {
			// Don't use .Parallel() without pinning
			t.Parallel()

			show, errGE := scrape.Episodes(s, el)
			is.NoErr(errGE)

			seasonNumbers, haveShowNumbers := showSeasonEpisodes[s]
			is.True(haveShowNumbers)

			for i := 0; i < len(show.Seasons); i++ {
				// If we have numbers for this season, make sure they match, then move on to the next season
				if episodeCount, haveSeasonNumbers := seasonNumbers[i+1]; haveSeasonNumbers {
					is.Equal(len(show.Seasons[i].Episodes), episodeCount)

					continue
				}

				// Otherwise, make sure we have at least one episode to inspect before going ahead
				lastEpIdx := len(show.Seasons[i].Episodes) - 1

				if lastEpIdx < 0 {
					continue
				}

				// Check the retrieved airdate against today's date
				if show.Seasons[i].Episodes[lastEpIdx].Airdate.Before(time.Now()) {
					t.Fatalf("missing episode count for S%02d of '%s' which has finished airing", i+1, show.Name)
				}
			}
		})
	}
}
