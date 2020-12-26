// Package collection defines an interface with which shows, seasons, and episodes can be stored and retrieved.
package collection

import (
	"go.jlucktay.dev/arrowverse/pkg/models"
)

type Shows interface {
	// InOrder will return the given show(s) in airdate order.
	InOrder(shows []string) ([]*models.Episode, error)

	// Count returns the number of shows in the collection.
	Count() (int, error)
	// Get returns the given show from the collection.
	Get(show string) (*models.Show, error)
	// GetAll returns all shows in the collection.
	GetAll() ([]*models.Show, error)

	// SeasonCount returns the number of seasons for the given show in the collection.
	SeasonCount(show string) (int, error)
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
