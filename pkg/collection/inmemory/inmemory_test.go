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
