package WebFramework

import "net/http"

type UrlPattern struct {
	UrlMapping map[string] interface{Handle(http.ResponseWriter, *http.Request)}
}

func (urlPattern *UrlPattern)AppendUrlPattern(uri string, v interface{Handle(http.ResponseWriter, *http.Request)}) {
	http.HandleFunc(uri, v.Handle)
}

func (urlPattern *UrlPattern)Init() {
	for key, value := range urlPattern.UrlMapping {
		urlPattern.AppendUrlPattern(key, value)
	}
}
