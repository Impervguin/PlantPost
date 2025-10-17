package urllib

import "net/url"

type URLStrategy interface {
	GetURL(path string) string
}

type StaticURLStrategy struct {
	BaseURL string
}

func (s *StaticURLStrategy) GetURL(path string) string {
	res, err := url.JoinPath(s.BaseURL, path)
	if err != nil {
		return ""
	}
	return res
}
