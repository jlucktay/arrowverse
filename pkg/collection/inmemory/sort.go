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

package inmemory

import (
	"fmt"
	"regexp"
	"strings"

	"go.jlucktay.dev/arrowverse/pkg/models"
)

type subexpressionName string

func (s subexpressionName) String() string {
	return string(s)
}

const (
	multiPartNumber subexpressionName = "mpn"
	multiPartTitle  subexpressionName = "mpt"
)

var (
	rePart = regexp.MustCompile(`^(?P<` + multiPartTitle.String() + `>.*) Part (?P<` + multiPartNumber.String() +
		`>[0-9a-zA-Z]+)$`)

	numbers = strings.NewReplacer("One", "1", "Two", "2", "Three", "3", "Four", "4", "Five", "5", "Six", "6",
		"Seven", "7", "Eight", "8", "Nine", "9")

	airOrderByYear = map[int][]models.ShowName{
		2016: {
			models.DCsLegendsOfTomorrow,
			models.Vixen,
		},

		2017: {
			models.TheFlashTheCW,
			models.DCsLegendsOfTomorrow,
		},

		2018: {
			models.TheFlashTheCW,
			models.Supergirl,
			models.BlackLightning,
			models.Arrow,
			models.DCsLegendsOfTomorrow,
		},

		2019: {
			models.DCsLegendsOfTomorrow,
			models.TheFlashTheCW,
			models.Arrow,
			models.BlackLightning,
			models.Batwoman,
			models.Supergirl,
		},

		2020: {
			models.Supergirl,
			models.Batwoman,
			models.TheFlashTheCW,
			models.Arrow,
			models.DCsLegendsOfTomorrow,
		},
	}
)

// ByAirdate implements sort.Interface for []*models.Episode based on the Airdate field.
type ByAirdate []models.Episode

func (a ByAirdate) Len() int      { return len(a) }
func (a ByAirdate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByAirdate) Less(i, j int) bool {
	// KISS; compare airdates first, and ONLY if they differ, go ahead with other logic
	if !a[i].Airdate.Equal(a[j].Airdate) {
		return a[i].Airdate.Before(a[j].Airdate)
	}

	// If the two episodes are from the same show, go by overall episode number
	if a[i].Season.Show.Name == a[j].Season.Show.Name {
		return a[i].EpisodeOverall < a[j].EpisodeOverall
	}

	// So now the two episodes aren't from the same show, but have the same airdate, and could be a multi-part/crossover
	if rePart.MatchString(a[i].Name) && rePart.MatchString(a[j].Name) {
		// At this point, both episode titles match the regex and contain '... Part ...', and are not from the same show
		iMatches := rePart.FindStringSubmatch(a[i].Name)
		jMatches := rePart.FindStringSubmatch(a[j].Name)

		titleIndex := rePart.SubexpIndex(multiPartTitle.String())

		// True if the multi-part/crossover has the same title, e.g. 'Crisis on Infinite Earths'
		if iMatches[titleIndex] == jMatches[titleIndex] {
			partNumberIndex := rePart.SubexpIndex(multiPartNumber.String())

			iPartNumber := numbers.Replace(iMatches[partNumberIndex])
			jPartNumber := numbers.Replace(jMatches[partNumberIndex])

			return iPartNumber < jPartNumber
		}
	}

	// If the previous block has not returned and flow has fallen through to this logic, at this point the two episodes
	// are not from the same show nor are they part of the same multi-part/crossover, but they did air on the same date
	// so we need to go through the air order, which is different in different years
	checkAirOrderYear := 0

	for year := range airOrderByYear {
		if year == a[i].Airdate.Year() {
			checkAirOrderYear = year

			break
		}
	}

	if checkAirOrderYear == 0 {
		return a[i].Airdate.Before(a[j].Airdate)
	}

	for _, showname := range airOrderByYear[checkAirOrderYear] {
		if showname == a[i].Season.Show.Name {
			return true
		}

		if showname == a[j].Season.Show.Name {
			return false
		}
	}

	// Neither show name turned up in the air orderings, so ... TODO?
	panic(fmt.Sprintf("Not sure what to do with these two episodes at this point:\n%s\n%s\n", a[i], a[j]))
}
