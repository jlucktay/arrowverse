package models

import (
	"fmt"
	"strings"
)

type ShowName string

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
