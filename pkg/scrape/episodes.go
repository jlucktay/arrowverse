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

package scrape

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/cenkalti/backoff/v4"
	"github.com/gocolly/colly/v2"

	"go.jlucktay.dev/arrowverse/pkg/models"
	"go.jlucktay.dev/arrowverse/pkg/util"
)

// Episodes will retrieve details for all of the given show's episodes from the wiki.
func Episodes(show models.ShowName, episodeListURL string) (*models.Show, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomain),
		colly.MaxDepth(0),
	)

	// Create the new show
	s := &models.Show{Name: show}

	c.OnHTML("body", func(body *colly.HTMLElement) {
		body.ForEach("table.wikitable", func(i int, table *colly.HTMLElement) {
			// Add a new season for this wikitable
			s.Seasons = append(s.Seasons, models.Season{Show: s, Number: i + 1})

			table.ForEach("tbody tr", func(rowNum int, tbody *colly.HTMLElement) {
				ep, errPE := processTableBody(tbody, rowNum, s, &s.Seasons[i])
				if errPE != nil {
					return
				}

				if ep != nil {
					// Add this episode to the current season, indexed by 'i' from body.ForEach
					s.Seasons[i].Episodes = append(s.Seasons[i].Episodes, *ep)
				}
			})
		})
	})

	// Execute the visit to actually make the HTTP request(s), inside an exponential backoff
	operation := func() error {
		var err error

		if errVisit := c.Visit(episodeListURL); errVisit != nil {
			err = fmt.Errorf("error visiting '%s': %w", episodeListURL, errVisit)
			fmt.Println(err)
		}

		return err
	}

	eb := backoff.NewExponentialBackOff()
	eb.MaxInterval = time.Second * 10
	eb.MaxElapsedTime = time.Minute

	if errVis := backoff.Retry(operation, eb); errVis != nil {
		return nil, fmt.Errorf("error while visiting %s: %w", episodeListURL, errVis)
	}

	return s, nil
}

func processTableBody(tbody *colly.HTMLElement, rowNum int, show *models.Show, season *models.Season) (*models.Episode, error) {
	if tbody.DOM.ChildrenFiltered("th").Length() > 0 { // Skip <th> row
		return nil, ErrNotEpisodeRow
	}

	var err error

	ep, itSel := &models.Episode{Season: season}, util.NewIteratingSelector()

	// Trim citation link suffixes like "[3]"
	checkCiteSuffix := regexp.MustCompile(`"?\[[0-9]+\]$`)

	if tbody.DOM.ChildrenFiltered("td").Length() >= 4 {
		raw := tbody.ChildText(itSel.Next())
		if ep.EpisodeOverall, err = strconv.Atoi(strings.TrimSpace(raw)); err != nil {
			return nil, fmt.Errorf("%q: %w", raw, ErrCouldNotParse)
		}
	} else {
		ep.EpisodeOverall = rowNum
	}

	epSeason := tbody.ChildText(itSel.Next())
	epSeason = checkCiteSuffix.ReplaceAllString(epSeason, "")

	// Handle the 'DC's Legends of Tomorrow' season 5 special episode
	if epSeason == `—` && ep.Season.Number == 5 && show.Name == models.DCsLegendsOfTomorrow {
		ep.EpisodeSeason = 0
	} else if ep.EpisodeSeason, err = strconv.Atoi(strings.TrimSpace(epSeason)); err != nil {
		return nil, fmt.Errorf("%q: %w", epSeason, ErrCouldNotParse)
	}

	epName := strings.Trim(strings.TrimSpace(tbody.ChildText(itSel.Next())), `"`)
	ep.Name = checkCiteSuffix.ReplaceAllString(epName, "")

	// Get ahead of too much junk data creeping in
	if ep.Name == "TBA" {
		return nil, ErrEpisodeNameTBA
	}

	rawURL := tbody.Request.AbsoluteURL(tbody.ChildAttr(itSel.String()+" a", "href"))
	if ep.Link, err = url.Parse(rawURL); err != nil {
		return nil, fmt.Errorf("%q: %w", rawURL, ErrCouldNotParse)
	}

	epAirdate := strings.TrimSpace(strings.Map(mapSpaces, tbody.ChildText(itSel.Next())))
	epAirdate = checkCiteSuffix.ReplaceAllString(epAirdate, "")

	// Round off 'TBA' airdates into the future 🤷‍♂️
	if epAirdate == "TBA" {
		theFuture := 5252 - time.Now().Year()

		ep.Airdate = time.Now().AddDate(theFuture, 0, 0).Round(time.Hour * 24)
	} else if ep.Airdate, err = time.Parse(models.AirdateLayout, epAirdate); err != nil {
		// First attempt to parse date/month/year from airdate, otherwise fall back to parsing just a year
		if ep.Airdate, err = time.Parse("2006", epAirdate); err != nil {
			return nil, fmt.Errorf("%q: %w", epAirdate, ErrCouldNotParse)
		}
	}

	return ep, nil
}

// mapSpaces helps us get rid of non-breaking spaces from HTML.
func mapSpaces(input rune) rune {
	if unicode.IsSpace(input) {
		return ' '
	}

	return input
}
