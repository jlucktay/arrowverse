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
func (c *Collection) AddSeason(show string, season *models.Season) error {
	panic("not implemented") // TODO: Implement
}

// AddEpisode to the given show's season in the collection.
func (c *Collection) AddEpisode(show string, season int, episode *models.Episode) error {
	panic("not implemented") // TODO: Implement
}

// InOrder will return episodes (limited to the given show(s) if any) in airdate order.
func (c *Collection) InOrder(shows ...string) ([]models.Episode, error) {
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
	panic("not implemented") // TODO: Implement
}

// Get returns the given show from the collection.
func (c *Collection) Get(show string) (*models.Show, error) {
	panic("not implemented") // TODO: Implement
}

// GetAll returns all shows in the collection.
func (c *Collection) GetAll() ([]*models.Show, error) {
	panic("not implemented") // TODO: Implement
}

// SeasonCount returns the number of seasons for the given show in the collection.
func (c *Collection) SeasonCount(show string) (int, error) {
	panic("not implemented") // TODO: Implement
}

// GetSeason returns the whole season for the given show in the collection.
func (c *Collection) GetSeason(show string, season int) (*models.Season, error) {
	panic("not implemented") // TODO: Implement
}

// GetAllSeasons returns all seasons in the collection for the given show.
func (c *Collection) GetAllSeasons(show string) ([]*models.Season, error) {
	panic("not implemented") // TODO: Implement
}

// EpisodeCount returns the number of episodes in the given show's season in the collection.
func (c *Collection) EpisodeCount(show string, season int) (int, error) {
	panic("not implemented") // TODO: Implement
}

// GetEpisode returns the specific episode from the given show's season in the collection.
func (c *Collection) GetEpisode(show string, season int, episode int) (*models.Episode, error) {
	panic("not implemented") // TODO: Implement
}

// GetAllSeasonEpisodes returns all episode from the given show's season in the collection.
func (c *Collection) GetAllSeasonEpisodes(show string, season int) ([]*models.Episode, error) {
	panic("not implemented") // TODO: Implement
}

// GetAllEpisodes returns all episodes from all seasons for the given show in the collection.
func (c *Collection) GetAllEpisodes(show string) ([]*models.Episode, error) {
	panic("not implemented") // TODO: Implement
}
