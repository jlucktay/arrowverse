package main

import (
	"fmt"
	"os"

	"go.jlucktay.dev/arrowverse/pkg/models"
	"go.jlucktay.dev/arrowverse/pkg/scrape"
)

func main() {
	episodeLists, errPE := scrape.EpisodeLists()
	if errPE != nil {
		fmt.Fprintf(os.Stderr, "could not get episode list URLs: %v", errPE)
	}

	shows := []models.Show{}

	for s, el := range episodeLists {
		show, errPE := scrape.Episodes(s, el)
		if errPE != nil {
			fmt.Fprintf(os.Stderr, "could not get episode details for '%s': %v", s, errPE)
		}

		shows = append(shows, *show)
	}

	for i := range shows {
		fmt.Println(shows[i])
	}
}
