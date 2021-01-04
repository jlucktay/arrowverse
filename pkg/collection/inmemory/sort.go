package inmemory

import (
	"regexp"
	"strings"

	"go.jlucktay.dev/arrowverse/pkg/models"
)

type subexpressionName string

func (s subexpressionName) String() string {
	return string(s)
}

const (
	multiPartNumber subexpressionName = "mpn"
	multiPartTitle  subexpressionName = "mpt"
)

var (
	rePart = regexp.MustCompile(`^(?P<` + multiPartTitle.String() + `>.*) Part (?P<` + multiPartNumber.String() +
		`>[0-9a-zA-Z]+)$`)

	numbers = strings.NewReplacer("One", "1", "Two", "2", "Three", "3", "Four", "4", "Five", "5", "Six", "6",
		"Seven", "7", "Eight", "8", "Nine", "9")
)

// ByAirdate implements sort.Interface for []*models.Episode based on the Airdate field.
type ByAirdate []models.Episode

func (a ByAirdate) Len() int      { return len(a) }
func (a ByAirdate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByAirdate) Less(i, j int) bool {
	// True if the two episodes aren't from the same show
	if a[i].Season.Show.Name != a[j].Season.Show.Name {
		if !rePart.MatchString(a[i].Name) || !rePart.MatchString(a[j].Name) {
			return a[i].Airdate.Before(a[j].Airdate)
		}

		// At this point, both episode titles match the regex and contain 'Part ...'

		iMatches := rePart.FindStringSubmatch(a[i].Name)
		jMatches := rePart.FindStringSubmatch(a[j].Name)

		titleIndex := rePart.SubexpIndex(multiPartTitle.String())

		// True if the multi-episode series has the same title, e.g. 'Crisis on Infinite Earths'
		if iMatches[titleIndex] == jMatches[titleIndex] {
			partNumberIndex := rePart.SubexpIndex(multiPartNumber.String())

			return numbers.Replace(iMatches[partNumberIndex]) < numbers.Replace(jMatches[partNumberIndex])
		}
	}

	// ...otherwise, fall through to remaining logic ... TODO?
	return a[i].EpisodeOverall < a[j].EpisodeOverall
}
