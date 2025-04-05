package postfilters

import (
	registry "PlantSite/internal/infra/postgres/filters/registry"
	"PlantSite/internal/models/search"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPostFilter(search.PostTagFilterID, PostTagsFilterFactory)
}

var _ registry.PostFilterFactory = PostTagsFilterFactory

func PostTagsFilterFactory(ps search.PostFilter) (registry.PostgresPostFilter, error) {
	pf, ok := ps.(*search.PostTagFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	tagSubquery := squirrel.Select("post_id").
		From("post_tag").
		Where(squirrel.Eq{"tag": pf.Tags}).
		GroupBy("post_id")

	filt := squirrel.Expr("id IN (?)", tagSubquery)

	return filt, nil
}
