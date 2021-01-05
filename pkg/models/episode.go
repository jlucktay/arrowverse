package models

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Episode describes an episode of an Arrowverse show.
type Episode struct {
	// Name of the episode.
	Name string

	// Season points back to the season containing this episode.
	Season *Season

	// EpisodeSeason is the episode number within the current season.
	EpisodeSeason int

	// EpisodeOverall is the episode number in the overall run of the entire show.
	EpisodeOverall int

	// Airdate is when the episode was first broadcast.
	Airdate time.Time

	// Link to a wiki page with episode details.
	Link *url.URL
}

func (e Episode) String() string {
	var b strings.Builder

	if e.Season != nil {
		if e.Season.Show != nil {
			fmt.Fprintf(&b, "%25s ", e.Season.Show.Name)
		} else {
			fmt.Fprintf(&b, "%25s ", "<unknown show>")
		}

		fmt.Fprintf(&b, "S%02d", e.Season.Number)
	} else {
		fmt.Fprint(&b, "S??")
	}

	fmt.Fprintf(&b, "E%02d %-61.61s", e.EpisodeSeason, e.Name)

	if !e.Airdate.IsZero() {
		fmt.Fprintf(&b, " %-20s", e.Airdate.Format(AirdateLayout))
	}

	if e.Link != nil {
		fmt.Fprintf(&b, " %s", e.Link)
	}

	return b.String()
}
