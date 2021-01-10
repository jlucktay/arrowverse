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

package inmemory_test

import (
	"errors"
	"testing"

	"github.com/matryer/is"

	"go.jlucktay.dev/arrowverse/pkg/collection/inmemory"
	"go.jlucktay.dev/arrowverse/pkg/models"
)

func TestAdd(t *testing.T) {
	im := inmemory.Collection{}
	is := is.New(t)

	c, err := im.Count()
	is.NoErr(err)
	is.Equal(c, 0)

	is.NoErr(im.Add(&models.Show{Name: models.Arrow}))

	c, err = im.Count()
	is.NoErr(err)
	is.Equal(c, 1)

	is.NoErr(im.Add(&models.Show{Name: models.Batwoman}))

	c, err = im.Count()
	is.NoErr(err)
	is.Equal(c, 2)
}

func TestAddSeason(t *testing.T) {
	im := inmemory.Collection{}
	is := is.New(t)

	_, err := im.SeasonCount(models.Arrow)
	is.True(errors.Is(err, inmemory.ErrNoSuchShow))

	err = im.AddSeason(models.Arrow, &models.Season{Number: 1})
	is.NoErr(err)

	c, err := im.SeasonCount(models.Arrow)
	is.NoErr(err)
	is.Equal(c, 1)
}
