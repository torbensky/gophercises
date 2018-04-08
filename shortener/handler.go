package shortener

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if val, ok := pathsToUrls[request.URL.Path]; ok {
			http.Redirect(response, request, val, http.StatusFound)
			return
		}
		fallback.ServeHTTP(response, request)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	mappings, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}

	handler := MapHandler(buildMap(mappings), fallback)

	return handler, nil
}

func buildMap(mappings []shortmap) map[string]string {
	r := make(map[string]string, len(mappings))
	for _, s := range mappings {
		r[s.Path] = s.URL
	}
	return r
}

func parseYAML(yml []byte) (shortmaps, error) {
	var mappings shortmaps
	err := yaml.Unmarshal(yml, &mappings)
	if err != nil {
		return nil, err
	}

	return mappings, nil
}

type shortmaps []shortmap
type shortmap struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
