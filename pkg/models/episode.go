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
