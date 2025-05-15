package urllib

import "net/url"

type UrlStrategy interface {
	GetUrl(path string) string
}

type StaticUrlStrategy struct {
	BaseUrl string
}

func (s *StaticUrlStrategy) GetUrl(path string) string {
	res, err := url.JoinPath(s.BaseUrl, path)
	if err != nil {
		return ""
	}
	return res
}
