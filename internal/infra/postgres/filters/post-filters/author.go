package postfilters

import (
	registry "PlantSite/internal/infra/postgres/filters/registry"
	"PlantSite/internal/models/search"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPostFilter(search.PostAuthorFilterID, PostAuthorFilterFactory)
}

var _ registry.PostFilterFactory = PostAuthorFilterFactory

func PostAuthorFilterFactory(ps search.PostFilter) (registry.PostgresPostFilter, error) {
	pf, ok := ps.(*search.PostAuthorFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	filt := squirrel.Eq{
		"author": pf.AuthorID,
	}

	return filt, nil
}
