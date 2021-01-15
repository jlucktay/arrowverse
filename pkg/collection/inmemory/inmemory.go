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

// Package inmemory implements in-memory storage to satisfy the Collection interface.
package inmemory

import (
	"sort"

	"go.jlucktay.dev/arrowverse/pkg/models"
)

type Collection struct {
	Shows []models.Show
}

// Add a show to the collection.
func (c *Collection) Add(s *models.Show) error {
	c.Shows = append(c.Shows, *s)

	return nil
}

// AddSeason to the given show in the collection.
func (c *Collection) AddSeason(show models.ShowName, season *models.Season) error {
	for i := range c.Shows {
		if c.Shows[i].Name == show {
			season.Show = &c.Shows[i]
			c.Shows[i].Seasons = append(c.Shows[i].Seasons, *season)

			return nil
		}
	}

	newShow := &models.Show{Name: show}
	season.Show = newShow
	newShow.Seasons = append(newShow.Seasons, *season)

	return c.Add(newShow)
}

// AddEpisode to the given show's season in the collection.
func (c *Collection) AddEpisode(show models.ShowName, season int, episode *models.Episode) error {
	for i := range c.Shows {
		if c.Shows[i].Name == show {
			for j := range c.Shows[i].Seasons {
				if c.Shows[i].Seasons[j].Number == season {
					episode.Season = &c.Shows[i].Seasons[j]
					c.Shows[i].Seasons[j].Episodes = append(c.Shows[i].Seasons[j].Episodes, *episode)

					return nil
				}
			}

			newSeason := models.Season{Show: &c.Shows[i], Number: season, Episodes: []models.Episode{*episode}}
			newSeason.Episodes[0].Season = &newSeason
			c.Shows[i].Seasons = append(c.Shows[i].Seasons, newSeason)

			return nil
		}
	}

	newShow := &models.Show{Name: show}
	newSeason := models.Season{Show: newShow, Number: season, Episodes: []models.Episode{*episode}}
	newSeason.Episodes[0].Season = &newSeason
	newShow.Seasons = append(newShow.Seasons, newSeason)

	return c.Add(newShow)
}

// InOrder will return episodes (limited to the given show(s) if any) in airdate order.
func (c *Collection) InOrder(shows ...models.ShowName) ([]models.Episode, error) {
	ret := []models.Episode{}

	for i := range c.Shows {
		if len(shows) > 0 {
			found := false

			for j := range shows {
				if shows[j] == c.Shows[i].Name {
					found = true

					break
				}
			}

			if !found {
				continue
			}
		}

		for k := range c.Shows[i].Seasons {
			ret = append(ret, c.Shows[i].Seasons[k].Episodes...)
		}
	}

	sort.Stable(ByAirdate(ret))

	return ret, nil
}

// Count returns the number of shows in the collection.
func (c *Collection) Count() (int, error) {
	return len(c.Shows), nil
}

// Get returns the given show from the collection.
func (c *Collection) Get(show string) (*models.Show, error) {
	panic("not implemented") // TODO: Implement.
}

// GetAll returns all shows in the collection.
func (c *Collection) GetAll() ([]*models.Show, error) {
	panic("not implemented") // TODO: Implement.
}

// SeasonCount returns the number of seasons for the given show in the collection.
func (c *Collection) SeasonCount(show models.ShowName) (int, error) {
	for i := range c.Shows {
		if c.Shows[i].Name == show {
			return len(c.Shows[i].Seasons), nil
		}
	}

	return -1, ErrNoSuchShow
}

// GetSeason returns the whole season for the given show in the collection.
func (c *Collection) GetSeason(show string, season int) (*models.Season, error) {
	panic("not implemented") // TODO: Implement.
}

// GetAllSeasons returns all seasons in the collection for the given show.
func (c *Collection) GetAllSeasons(show string) ([]*models.Season, error) {
	panic("not implemented") // TODO: Implement.
}

// EpisodeCount returns the number of episodes in the given show's season in the collection.
func (c *Collection) EpisodeCount(show string, season int) (int, error) {
	panic("not implemented") // TODO: Implement.
}

// GetEpisode returns the specific episode from the given show's season in the collection.
func (c *Collection) GetEpisode(show string, season, episode int) (*models.Episode, error) {
	panic("not implemented") // TODO: Implement.
}

// GetAllSeasonEpisodes returns all episode from the given show's season in the collection.
func (c *Collection) GetAllSeasonEpisodes(show string, season int) ([]*models.Episode, error) {
	panic("not implemented") // TODO: Implement.
}

// GetAllEpisodes returns all episodes from all seasons for the given show in the collection.
func (c *Collection) GetAllEpisodes(show string) ([]*models.Episode, error) {
	panic("not implemented") // TODO: Implement.
}
