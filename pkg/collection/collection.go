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

// Package collection defines an interface with which shows, seasons, and episodes can be stored and retrieved.
package collection

import (
	"go.jlucktay.dev/arrowverse/pkg/models"
)

type Shows interface {
	// Add a show to the collection.
	Add(*models.Show) error
	// AddSeason to the given show in the collection.
	AddSeason(show models.ShowName, season *models.Season) error
	// AddEpisode to the given show's season in the collection.
	AddEpisode(show models.ShowName, season int, episode *models.Episode) error

	// InOrder will return episodes (limited to the given show(s) if any) in airdate order.
	InOrder(shows ...models.ShowName) ([]models.Episode, error)

	// Count returns the number of shows in the collection.
	Count() (int, error)
	// Get returns the given show from the collection.
	Get(show string) (*models.Show, error)
	// GetAll returns all shows in the collection.
	GetAll() ([]*models.Show, error)

	// SeasonCount returns the number of seasons for the given show in the collection.
	SeasonCount(show models.ShowName) (int, error)
	// GetSeason returns the whole season for the given show in the collection.
	GetSeason(show string, season int) (*models.Season, error)
	// GetAllSeasons returns all seasons in the collection for the given show.
	GetAllSeasons(show string) ([]*models.Season, error)

	// EpisodeCount returns the number of episodes in the given show's season in the collection.
	EpisodeCount(show string, season int) (int, error)
	// GetEpisode returns the specific episode from the given show's season in the collection.
	GetEpisode(show string, season, episode int) (*models.Episode, error)
	// GetAllSeasonEpisodes returns all episode from the given show's season in the collection.
	GetAllSeasonEpisodes(show string, season int) ([]*models.Episode, error)
	// GetAllEpisodes returns all episodes from all seasons for the given show in the collection.
	GetAllEpisodes(show string) ([]*models.Episode, error)
}
