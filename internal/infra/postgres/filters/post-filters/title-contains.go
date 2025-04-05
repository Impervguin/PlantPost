package postfilters

import (
	registry "PlantSite/internal/infra/postgres/filters/registry"
	"PlantSite/internal/models/search"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPostFilter(search.PostTitleContainsFilterID, PostTitleContainsFilterFactory)
}

var _ registry.PostFilterFactory = PostTitleContainsFilterFactory

func PostTitleContainsFilterFactory(ps search.PostFilter) (registry.PostgresPostFilter, error) {
	pf, ok := ps.(*search.PostTitleContainsFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	filt := squirrel.ILike{
		"title": fmt.Sprintf("%%%s%%", pf.Part),
	}

	return filt, nil
}
