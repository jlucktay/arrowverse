package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"go.jlucktay.dev/arrowverse/pkg/collection"
	"go.jlucktay.dev/arrowverse/pkg/collection/inmemory"
	"go.jlucktay.dev/arrowverse/pkg/models"
	"go.jlucktay.dev/arrowverse/pkg/scrape"
)

func runRoot(cmd *cobra.Command, _ []string) error {
	episodeLists, errEL := scrape.EpisodeLists()
	if errEL != nil {
		return fmt.Errorf("could not get episode lists: %w", errEL)
	}

	var cs collection.Shows = &inmemory.Collection{}

	for s, el := range episodeLists {
		show, errEps := scrape.Episodes(s, el)
		if errEps != nil {
			return fmt.Errorf("could not get episode details for '%s': %w", s, errEps)
		}

		if errAdd := cs.Add(show); errAdd != nil {
			return fmt.Errorf("could not add '%s' details to collection: %w", show.Name, errAdd)
		}
	}

	includedShows := []models.ShowName{
		models.Arrow,
		models.Batwoman,
		models.BlackLightning,
		models.Constantine,
		models.DCsLegendsOfTomorrow,
		models.FreedomFightersTheRay,
		models.Supergirl,
		models.TheFlashTheCW,
		models.Vixen,

		// // These two Arrowverse shows aren't included by default
		// models.BirdsOfPrey,
		// models.TheFlashCBS,
	}

	csio, errIO := cs.InOrder(includedShows...)
	if errIO != nil {
		return fmt.Errorf("could not get episode details: %w", errIO)
	}

	for i := range csio {
		if csio[i].Airdate.Year() < 5252 { //nolint:gomnd // https://dc.fandom.com/wiki/52#52
			fmt.Fprintf(cmd.OutOrStdout(), "[%03d] %s\n", i, csio[i])
		}
	}

	return nil
}
