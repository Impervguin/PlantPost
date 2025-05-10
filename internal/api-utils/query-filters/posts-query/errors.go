package postsquery

import "errors"

var (
	ErrParserNotFound = errors.New("no parser registered for this filter")
	ErrParsingFailed  = errors.New("parser failed to parse filter")
)
