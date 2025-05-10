package postsquery

import (
	"PlantSite/internal/models/search"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func parsePostTitleFilterfunc(queryValue string) (search.PostFilter, error) {
	filt := search.NewPostTitleContainsFilter(queryValue)
	if filt == nil {
		return nil, ErrParsingFailed
	}
	return filt, nil
}

func parsePostTagsFilterfunc(queryValue string) (search.PostFilter, error) {
	// var1,var2,... format
	tags := strings.Split(queryValue, ",")
	if len(tags) == 0 {
		return nil, ErrParsingFailed
	}
	filt := search.NewPostTagFilter(tags)
	if filt == nil {
		return nil, ErrParsingFailed
	}
	return filt, nil
}

func parsePostAuthorFilterfunc(queryValue string) (search.PostFilter, error) {
	authorID, err := uuid.Parse(queryValue)
	if err != nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PostAuthorFilterParam, queryValue)
	}
	filt := search.NewPostAuthorFilter(authorID)
	if filt == nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PostAuthorFilterParam, queryValue)
	}
	return filt, nil
}
