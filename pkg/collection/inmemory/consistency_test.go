package inmemory_test

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/google/go-cmp/cmp"

	"go.jlucktay.dev/arrowverse/pkg/collection"
	"go.jlucktay.dev/arrowverse/pkg/collection/inmemory"
	"go.jlucktay.dev/arrowverse/pkg/models"
	"go.jlucktay.dev/arrowverse/pkg/scrape"
	"go.jlucktay.dev/arrowverse/pkg/util"
)

func TestConsistencyWithArrowverseDotInfo(t *testing.T) {
	t.SkipNow()

	// get episodes from arrow.fandom.com
	// get episodes from arrowverse.info
	// compare with go-cmp

	const (
		host    = "arrowverse.info"
		fullURL = "https://" + host
	)

	reEpisode := regexp.MustCompile(`^S(?P<season>[0-9]{2})E(?P<episode>[0-9]{2})$`)

	c := colly.NewCollector(
		colly.AllowedDomains(host),
		colly.MaxDepth(0),
		colly.Debugger(&debug.LogDebugger{}),
	)

	// Create somewhere to store the shows
	var csArrowverseInfo collection.Shows = &inmemory.Collection{}

	c.OnHTML("body", func(body *colly.HTMLElement) {
		body.ForEach("table#episode-list", func(_ int, table *colly.HTMLElement) {
			table.ForEach("tbody tr.episode", func(_ int, tbody *colly.HTMLElement) {
				var (
					err error
					ep  models.Episode
				)

				itSel := util.NewIteratingSelector()

				overall := tbody.ChildText(itSel.Next())

				showNameCandidate := tbody.ChildText(itSel.Next())
				if !models.ValidShowName(showNameCandidate) {
					t.Fatalf("show name '%s' is not valid", showNameCandidate)
				}

				showName := models.ShowName(showNameCandidate)

				seasonAndEpisode := strings.Map(fixScraped, tbody.ChildText(itSel.Next()))
				ep.Name = tbody.ChildText(itSel.Next())
				airdate := tbody.ChildText(itSel.Next())

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

				seasonNumber, errConvSeason := strconv.Atoi(matches[seasonIndex])
				if errConvSeason != nil {
					t.Fatalf("could not parse season number '%s': %v", matches[seasonIndex], errConvSeason)
				}

				ep.Season = &models.Season{
					Show: &models.Show{
						Name: showName,
					},
					Number: seasonNumber,
				}
				t.Logf("[%3s] %s", overall, ep)
				ep.Season = nil

				if err = csArrowverseInfo.AddEpisode(showName, seasonNumber, &ep); err != nil {
					t.Fatalf("could not add episode '%#v' to collection: %v", ep, err)
				}
			})
		})
	})

	// Execute the visit to actually make the HTTP request(s), inside an exponential backoff with default settings
	operation := func() error {
		return c.Visit(fullURL)
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

	arrowverseInfoEpisodes, err := csArrowverseInfo.InOrder(includedShows...)
	if err != nil {
		t.Fatalf("error getting episodes in order: %v", err)
	}

	episodeLists, errEL := scrape.EpisodeLists()
	if errEL != nil {
		t.Fatalf("could not get episode lists: %v", errEL)
	}

	var csArrowFandomCom collection.Shows = &inmemory.Collection{}

	for s, el := range episodeLists {
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

	t.Logf("arrowverse.info #eps: %d", len(arrowverseInfoEpisodes))
	t.Logf("arrow.fandom.com #eps: %d", len(arrowFandomComEpisodes))

	t.Skip("TODO: finish this comparison") // TODO

	if diff := cmp.Diff(arrowverseInfoEpisodes, arrowFandomComEpisodes); diff != "" {
		t.Errorf("Mismatch (-arrowverse.info +arrow.fandom.com):\n%s", diff)
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
