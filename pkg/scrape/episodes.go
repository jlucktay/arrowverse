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

// GetEpisodes will retrieve details for all of the given show's episodes from the wiki.
func Episodes(show, episodeListURL string) (*models.Show, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomain),
		colly.MaxDepth(0),
	)

	// Create the new show
	s := &models.Show{Name: show}

	c.OnHTML("body", func(body *colly.HTMLElement) {
		body.ForEach("table.wikitable", func(i int, table *colly.HTMLElement) {
			// Add a new season for this wikitable
			s.Seasons = append(s.Seasons, models.Season{Number: i + 1})

			table.ForEach("tbody tr", func(_ int, tbody *colly.HTMLElement) {
				if tbody.DOM.ChildrenFiltered("th").Length() > 0 { // Skip <th> row
					return
				}

				var err error
				ep, itSel := models.Episode{}, util.NewIteratingSelector()

				// Trim citation link suffixes like "[3]"
				checkCiteSuffix := regexp.MustCompile(`"?\[[0-9]+\]$`)

				if tbody.DOM.ChildrenFiltered("td").Length() >= 4 { //nolint:gomnd // Deal with wider tables
					if ep.EpisodeOverall, err = strconv.Atoi(strings.TrimSpace(tbody.ChildText(itSel.Next()))); err != nil {
						return
					}
				}

				epSeason := tbody.ChildText(itSel.Next())
				epSeason = checkCiteSuffix.ReplaceAllString(epSeason, "")

				// Handle the 'DC's Legends of Tomorrow' season 5 special episode
				if epSeason == `—` && s.Seasons[i].Number == 5 && s.Name == "DC's Legends of Tomorrow" {
					ep.EpisodeSeason = 0
				} else {
					if ep.EpisodeSeason, err = strconv.Atoi(strings.TrimSpace(epSeason)); err != nil {
						return
					}
				}

				epName := strings.Trim(strings.TrimSpace(tbody.ChildText(itSel.Next())), `"`)
				ep.Name = checkCiteSuffix.ReplaceAllString(epName, "")

				// Get ahead of too much junk data creeping in
				if ep.Name == "TBA" {
					return
				}

				raw := tbody.Request.AbsoluteURL(tbody.ChildAttr(itSel.String()+" a", "href"))
				if ep.Link, err = url.Parse(raw); err != nil {
					return
				}

				epAirdate := strings.TrimSpace(strings.Map(mapSpaces, tbody.ChildText(itSel.Next())))
				epAirdate = checkCiteSuffix.ReplaceAllString(epAirdate, "")

				// Round off 'TBA' airdates into the future 🤷‍♂️
				if epAirdate == "TBA" {
					theFuture := 5252 - time.Now().Year() //nolint:gomnd // https://dc.fandom.com/wiki/52#52

					ep.Airdate = time.Now().AddDate(theFuture, 0, 0).Round(time.Hour * 24) //nolint:gomnd // 24h = 1d
				} else {
					if ep.Airdate, err = time.Parse(models.AirdateLayout, epAirdate); err != nil {
						return
					}
				}

				// Add this episode to the current season, indexed by 'i' from body.ForEach
				s.Seasons[i].Episodes = append(s.Seasons[i].Episodes, ep)
			})
		})
	})

	// Execute the visit to actually make the HTTP request(s), inside an exponential backoff with default settings
	operation := func() error {
		return c.Visit(episodeListURL)
	}

	if errVis := backoff.Retry(operation, backoff.NewExponentialBackOff()); errVis != nil {
		return nil, fmt.Errorf("error while visiting %s: %w", episodeListURL, errVis)
	}

	return s, nil
}

// mapSpaces helps us get rid of non-breaking spaces from HTML.
func mapSpaces(input rune) rune {
	if unicode.IsSpace(input) {
		return ' '
	}

	return input
}
