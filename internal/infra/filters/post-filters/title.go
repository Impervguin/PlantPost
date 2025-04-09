package postfilters

import (
	registry "PlantSite/internal/infra/filters/registry"
	"PlantSite/internal/models/search"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPostFilter(search.PostTitleFilterID, PostNameFilterFactory)
}

var _ registry.PostFilterFactory = PostNameFilterFactory

func PostNameFilterFactory(ps search.PostFilter) (registry.PostgresPostFilter, error) {
	pf, ok := ps.(*search.PostTitleFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	filt := squirrel.Eq{
		"title": pf.Title,
	}

	return filt, nil
}
