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
