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

package models

import (
	"fmt"
	"strings"
)

// ShowName is a custom string type for use with these models.
type ShowName string

// These are the names of the shows that this app wrangles.
const (
	Arrow                 ShowName = "Arrow"
	Batwoman              ShowName = "Batwoman"
	BlackLightning        ShowName = "Black Lightning"
	Constantine           ShowName = "Constantine"
	DCsLegendsOfTomorrow  ShowName = "DC's Legends of Tomorrow"
	FreedomFightersTheRay ShowName = "Freedom Fighters: The Ray"
	Supergirl             ShowName = "Supergirl"
	TheFlashTheCW         ShowName = "The Flash (The CW)"
	Vixen                 ShowName = "Vixen"

	BirdsOfPrey ShowName = "Birds of Prey"
	TheFlashCBS ShowName = "The Flash (CBS)"
)

// ValidShowName will assert whether or not the given string is a valid Show name.
func ValidShowName(s string) bool {
	validShowNames := []ShowName{
		Arrow,
		Batwoman,
		BirdsOfPrey,
		BlackLightning,
		Constantine,
		DCsLegendsOfTomorrow,
		FreedomFightersTheRay,
		Supergirl,
		TheFlashCBS,
		TheFlashTheCW,
		Vixen,
	}

	for i := range validShowNames {
		if strings.EqualFold(string(validShowNames[i]), s) {
			return true
		}
	}

	return false
}

// Show describes an Arrowverse show.
type Show struct {
	// Name of the show.
	Name ShowName

	// Seasons for this show only.
	Seasons []Season
}

func (s Show) String() string {
	var b strings.Builder

	for _, season := range s.Seasons {
		fmt.Fprintf(&b, "%s, season %d/%d (%d episode(s))\n",
			s.Name, season.Number, len(s.Seasons), len(season.Episodes))
		fmt.Fprintf(&b, "%s\n", season)
	}

	return b.String()
}
