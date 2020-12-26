package scrape

import "errors"

var (
	ErrNotEpisodeRow  = errors.New("not an episode row")
	ErrCouldNotParse  = errors.New("could not parse")
	ErrEpisodeNameTBA = errors.New("episode name has yet to be announced")
)
