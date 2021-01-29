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

package scrape

import (
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gocolly/colly/v2"

	"go.jlucktay.dev/arrowverse/pkg/models"
)

// EpisodeLists will retrieve URLs of all of the 'List of <show> episodes' for shows that are available on the wiki.
func EpisodeLists() (map[models.ShowName]string, error) {
	const (
		checkPrefix = "List of "
		checkSuffix = " episodes"
	)

	episodeListURLs := map[models.ShowName]string{}

	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomain),
		colly.MaxDepth(0),
	)

	c.OnHTML("body", func(body *colly.HTMLElement) {
		body.ForEach("div.category-page__members "+
			"ul.category-page__members-for-char "+
			"li.category-page__member "+
			"a.category-page__member-link",
			func(_ int, a *colly.HTMLElement) {
				// Only consider 'List of <show> episodes' links
				if !strings.HasPrefix(a.Text, checkPrefix) || !strings.HasSuffix(a.Text, checkSuffix) {
					return
				}

				// Need to specifically exclude the list of crossover episodes
				if strings.Contains(a.Text, " crossover ") {
					return
				}

				showName := strings.TrimPrefix(a.Text, checkPrefix)
				showName = strings.TrimSuffix(showName, checkSuffix)

				if !models.ValidShowName(showName) {
					return
				}

				episodeListURLs[models.ShowName(showName)] = a.Request.AbsoluteURL(a.Attr("href"))
			})
	})

	// Execute the visit to actually make the HTTP request(s), inside an exponential backoff
	operation := func() error {
		var err error

		if errVisit := c.Visit(categoryListsURL); errVisit != nil {
			err = fmt.Errorf("error visiting '%s': %w", categoryListsURL, errVisit)
			fmt.Println(err)
		}

		return err
	}

	eb := backoff.NewExponentialBackOff()
	eb.MaxInterval = time.Second * 10 //nolint:gomnd // Wait a maximum of 10 seconds between visit attempts.
	eb.MaxElapsedTime = time.Minute

	if errVis := backoff.Retry(operation, eb); errVis != nil {
		return nil, fmt.Errorf("error while visiting %s: %w", categoryListsURL, errVis)
	}

	return episodeListURLs, nil
}
