package postsquery

import (
	"PlantSite/internal/models/search"
	"sync"
)

type QueryPostFilterParser func(string) (search.PostFilter, error)

type PostFilterParam string

const (
	PostTitleFilterParam  PostFilterParam = "title"
	PostTagsFilterParam   PostFilterParam = "tags"
	PostAuthorFilterParam PostFilterParam = "author"
)

type QueryPostFilterRegistry struct {
	parsers map[PostFilterParam]QueryPostFilterParser
	mut     sync.RWMutex
}

func newQueryPostFilterRegistry() *QueryPostFilterRegistry {
	return &QueryPostFilterRegistry{
		parsers: make(map[PostFilterParam]QueryPostFilterParser),
		mut:     sync.RWMutex{},
	}
}

func (r *QueryPostFilterRegistry) register(name PostFilterParam, parser QueryPostFilterParser) {
	r.mut.Lock()
	defer r.mut.Unlock()
	r.parsers[name] = parser
}

func (r *QueryPostFilterRegistry) parse(name PostFilterParam, queryValue string) (search.PostFilter, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()
	parser, ok := r.parsers[name]
	if !ok {
		return nil, ErrParserNotFound
	}
	return parser(queryValue)
}
