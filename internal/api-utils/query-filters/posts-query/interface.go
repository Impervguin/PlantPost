package postsquery

import (
	"PlantSite/internal/models/search"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	registry = newQueryPostFilterRegistry()
	once     = sync.Once{}
)

func registryInit() {
	once.Do(func() {
		registry.register(PostTitleFilterParam, parsePostTitleFilterfunc)
		registry.register(PostTagsFilterParam, parsePostTagsFilterfunc)
		registry.register(PostAuthorFilterParam, parsePostAuthorFilterfunc)
	})
}

func ParseQueryPostFilter(filterType PostFilterParam, queryValue string) (search.PostFilter, error) {
	registryInit()
	return registry.parse(filterType, queryValue)
}

func ParseGinQueryPostSearch(c *gin.Context) (*search.PostSearch, error) {
	registryInit()
	params := c.Request.URL.Query()
	srch := search.NewPostSearch()
	for filterType, query := range params {
		if len(query) == 0 {
			continue
		}
		for _, q := range query {
			filter, err := registry.parse(PostFilterParam(filterType), q)
			if err != nil {
				return &search.PostSearch{}, fmt.Errorf("can't parse filter in search: %w", err)
			}
			srch.AddFilter(filter)
		}
	}
	return srch, nil
}
