package inmemory

import "go.jlucktay.dev/arrowverse/pkg/models"

// ByAirdate implements sort.Interface for []*models.Episode based on the Airdate field.
type ByAirdate []models.Episode

func (a ByAirdate) Len() int      { return len(a) }
func (a ByAirdate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByAirdate) Less(i, j int) bool {
	// If the two episodes aren't from the same show, then comparison of airdate is sufficient, for now. TODO?
	if a[i].Season.Show.Name != a[j].Season.Show.Name {
		return a[i].Airdate.Before(a[j].Airdate)
	}

	// Otherwise ... TODO?
	return a[i].EpisodeOverall < a[j].EpisodeOverall
}
