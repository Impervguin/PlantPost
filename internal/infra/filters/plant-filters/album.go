package plantfilters

import (
	"PlantSite/internal/models/search"

	registry "PlantSite/internal/infra/filters/registry"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.PlantAlbumFilterID, PlantAlbumFilterFactory)
}

var _ registry.PlantFilterFactory = PlantAlbumFilterFactory

func PlantAlbumFilterFactory(ps search.PlantFilter) (registry.PostgresPlantFilter, error) {
	pf, ok := ps.(*search.PlantAlbumFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	tagSubquery := squirrel.Select("plant_id").
		From("plant_album").
		Where(squirrel.Eq{"album_id": pf.AlbumID})

	filt := squirrel.Expr("id IN (?)", tagSubquery)

	return filt, nil
}
