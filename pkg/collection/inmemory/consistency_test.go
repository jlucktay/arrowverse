package inmemory_test

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gocolly/colly/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"go.jlucktay.dev/arrowverse/pkg/collection"
	"go.jlucktay.dev/arrowverse/pkg/collection/inmemory"
	"go.jlucktay.dev/arrowverse/pkg/models"
	"go.jlucktay.dev/arrowverse/pkg/scrape"
	"go.jlucktay.dev/arrowverse/pkg/util"
)

func TestConsistencyWithArrowverseDotInfo(t *testing.T) {
	const (
		host    = "arrowverse.info"
		fullURL = "https://" + host
	)

	reEpisode := regexp.MustCompile(`^S(?P<season>[0-9]{2})E(?P<episode>[0-9]{2})$`)

	c := colly.NewCollector(
		colly.AllowedDomains(host),
		colly.MaxDepth(0),
	)

	// Create somewhere to store the list from arrowverse.info, and keep track of per-show overall episode number
	var arrowverseInfoEpisodes []models.Episode
	episodeOverallCounters := map[models.ShowName]int{}

	c.OnHTML("body", func(body *colly.HTMLElement) {
		body.ForEach("table#episode-list", func(_ int, table *colly.HTMLElement) {
			table.ForEach("tbody tr.episode", func(_ int, tbody *colly.HTMLElement) {
				var (
					err error
					ep  models.Episode
				)

				itSel := util.NewIteratingSelector()

				_ = strings.TrimSpace(tbody.ChildText(itSel.Next()))

				showNameCandidate := strings.TrimSpace(tbody.ChildText(itSel.Next()))
				if showNameCandidate == "The Flash" {
					showNameCandidate += " (The CW)"
				}

				if !models.ValidShowName(showNameCandidate) {
					t.Fatalf("show name '%s' is not valid", showNameCandidate)
				}

				showName := models.ShowName(showNameCandidate)

				seasonAndEpisode := strings.Map(fixScraped, strings.TrimSpace(tbody.ChildText(itSel.Next())))
				ep.Name = strings.TrimSpace(tbody.ChildText(itSel.Next()))
				airdate := strings.TrimSpace(tbody.ChildText(itSel.Next()))

				if !reEpisode.MatchString(seasonAndEpisode) {
					t.Fatalf("regex could not parse season/episode from third column: %s", seasonAndEpisode)
				}

				matches := reEpisode.FindStringSubmatch(seasonAndEpisode)
				seasonIndex := reEpisode.SubexpIndex("season")
				episodeSeasonIndex := reEpisode.SubexpIndex("episode")

				if ep.EpisodeSeason, err = strconv.Atoi(matches[episodeSeasonIndex]); err != nil {
					t.Fatalf("could not convert season/episode '%s': %v", matches[episodeSeasonIndex], err)
				}

				if ep.Airdate, err = time.Parse(models.AirdateLayout, airdate); err != nil {
					t.Fatalf("could not parse airdate '%s': %v", airdate, err)
				}

				if ep.Airdate.Year() >= 5252 {
					return
				}

				if ep.Link, err = url.Parse(tbody.Attr("data-href")); err != nil {
					t.Fatalf("could not parse episode link '%s': %v", tbody.Attr("data-href"), err)
				}

				// https://www.eff.org/https-everywhere
				if ep.Link.Scheme == "http" {
					ep.Link.Scheme = "https"
				}

				seasonNumber, errConvSeason := strconv.Atoi(matches[seasonIndex])
				if errConvSeason != nil {
					t.Fatalf("could not parse season number '%s': %v", matches[seasonIndex], errConvSeason)
				}

				episodeOverallCounters[showName] += 1
				ep.EpisodeOverall = episodeOverallCounters[showName]

				ep.Season = &models.Season{
					Show:   &models.Show{Name: showName},
					Number: seasonNumber}

				// Workaround for Black Lightning pedantry
				if ep.Season.Show.Name == models.BlackLightning && ep.Season.Number == 3 &&
					((ep.EpisodeSeason >= 2 && ep.EpisodeSeason <= 7) || (ep.EpisodeSeason >= 11 && ep.EpisodeSeason <= 13) ||
						ep.EpisodeSeason == 15) && strings.HasPrefix(ep.Name, "The Book of ") {
					lastColon := strings.LastIndex(ep.Name, ":")

					if lastColon > -1 {
						ep.Name = ep.Name[0:lastColon]
					}
				}
				// Workaround for Black Lightning pedantry

				arrowverseInfoEpisodes = append(arrowverseInfoEpisodes, ep)
			})
		})
	})

	// Execute the visit to actually make the HTTP request(s), inside an exponential backoff with default settings
	operation := func() error {
		var err error

		if err = c.Visit(fullURL); err != nil {
			fmt.Printf("error visiting: %v", err)
		}

		return err
	}

	if errVis := backoff.Retry(operation, backoff.NewExponentialBackOff()); errVis != nil {
		t.Fatalf("error while visiting %s: %v", fullURL, errVis)
	}

	includedShows := []models.ShowName{
		models.Arrow,
		models.Batwoman,
		// models.BirdsOfPrey,
		models.BlackLightning,
		models.Constantine,
		models.DCsLegendsOfTomorrow,
		models.FreedomFightersTheRay,
		models.Supergirl,
		// models.TheFlashCBS,
		models.TheFlashTheCW,
		models.Vixen,
	}

	arrowFandomComEpisodeLists, errEL := scrape.EpisodeLists()
	if errEL != nil {
		t.Fatalf("could not get episode lists: %v", errEL)
	}

	var csArrowFandomCom collection.Shows = &inmemory.Collection{}

	for s, el := range arrowFandomComEpisodeLists {
		show, errEps := scrape.Episodes(s, el)
		if errEps != nil {
			t.Fatalf("could not get episode details for '%s': %v", s, errEps)
		}

		if errAdd := csArrowFandomCom.Add(show); errAdd != nil {
			t.Fatalf("could not add '%s' details to collection: %v", show.Name, errAdd)
		}
	}

	arrowFandomComEpisodes, errIO := csArrowFandomCom.InOrder(includedShows...)
	if errIO != nil {
		t.Fatalf("could not get episode details: %v", errIO)
	}

	t.Logf("arrowverse.info  #eps: %d", len(arrowverseInfoEpisodes))
	t.Logf("arrow.fandom.com #eps: %d", len(arrowFandomComEpisodes))

	ignoreSeasonAndLink := cmpopts.IgnoreFields(models.Episode{}, "Season", "Link")

	for i := 0; i < len(arrowverseInfoEpisodes) && i < len(arrowFandomComEpisodes) &&
		arrowverseInfoEpisodes[i].Airdate.Before(time.Now()) && arrowFandomComEpisodes[i].Airdate.Before(time.Now()); i++ {
		if diff := cmp.Diff(arrowverseInfoEpisodes[i], arrowFandomComEpisodes[i], ignoreSeasonAndLink); diff != "" {
			t.Logf("\narrowverse.info:  (%s) '%#v'\narrow.fandom.com: (%s) '%#v'\n%[2]s\n%[4]s",
				arrowverseInfoEpisodes[i].Season.Show.Name, arrowverseInfoEpisodes[i],
				arrowFandomComEpisodes[i].Season.Show.Name, arrowFandomComEpisodes[i])
			t.Fatalf("Mismatch (-arrowverse.info[%03[1]d] +arrow.fandom.com[%03[1]d]):\n%[2]s", i, diff)
		}
	}
}

// fixScraped helps us straighten out scraped data.
func fixScraped(input rune) rune {
	if input == '—' {
		return '0'
	}

	if input == '[' || input == ']' {
		return -1
	}

	return input
}
