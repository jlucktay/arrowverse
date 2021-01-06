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

	airOrderByYear = map[int][]models.ShowName{
		2016: {
			models.DCsLegendsOfTomorrow,
			models.Vixen,
		},

		2017: {
			models.TheFlashTheCW,
			models.DCsLegendsOfTomorrow,
		},

		2018: {
			models.TheFlashTheCW,
			models.Supergirl,
			models.BlackLightning,
			models.Arrow,
			models.DCsLegendsOfTomorrow,
		},

		2019: {
			models.DCsLegendsOfTomorrow,
			models.TheFlashTheCW,
			models.Arrow,
			models.BlackLightning,
			models.Batwoman,
			models.Supergirl,
		},
	}
)

// ByAirdate implements sort.Interface for []*models.Episode based on the Airdate field.
type ByAirdate []models.Episode

func (a ByAirdate) Len() int      { return len(a) }
func (a ByAirdate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByAirdate) Less(i, j int) bool {
	// If the two episodes are from the same show, go by airdate or overall episode number
	if a[i].Season.Show.Name == a[j].Season.Show.Name {
		return a[i].Airdate.Before(a[j].Airdate) || a[i].EpisodeOverall < a[j].EpisodeOverall
	}

	// So now the two episodes aren't from the same show, but could be multi-part
	if !rePart.MatchString(a[i].Name) || !rePart.MatchString(a[j].Name) {
		// At this point, they're not from the same show nor are they multi-part, so we check the year they both aired
		if a[i].Airdate.Year() != a[j].Airdate.Year() {
			return a[i].Airdate.Before(a[j].Airdate)
		}

		// They aired in the same year, so we may need to check the air order which is different in different years
		checkAirOrderYear := 0
		for year := range airOrderByYear {
			if year == a[i].Airdate.Year() && year == a[j].Airdate.Year() {
				checkAirOrderYear = year
				break
			}
		}

		if checkAirOrderYear == 0 {
			return a[i].Airdate.Before(a[j].Airdate)
		}

		for _, showname := range airOrderByYear[checkAirOrderYear] {
			if a[i].Season.Show.Name == showname {
				return true
			}
			if a[j].Season.Show.Name == showname {
				return false
			}
		}

		// Neither show name turned up in the air orderings, so fall back to airdate
		return a[i].Airdate.Before(a[j].Airdate)
	}

	// At this point, both episode titles match the regex and contain '... Part ...', and are not from the same show
	iMatches := rePart.FindStringSubmatch(a[i].Name)
	jMatches := rePart.FindStringSubmatch(a[j].Name)

	titleIndex := rePart.SubexpIndex(multiPartTitle.String())

	// True if the multi-episode series has the same title, e.g. 'Crisis on Infinite Earths'
	if iMatches[titleIndex] == jMatches[titleIndex] {
		partNumberIndex := rePart.SubexpIndex(multiPartNumber.String())

		return numbers.Replace(iMatches[partNumberIndex]) < numbers.Replace(jMatches[partNumberIndex])
	}

	// Otherwise, fall back on remaining logic
	return a[i].Airdate.Before(a[j].Airdate) || a[i].EpisodeOverall < a[j].EpisodeOverall
}
