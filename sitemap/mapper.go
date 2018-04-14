package sitemap

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"

	"github.com/torbensky/gophercises/links"
)

// SiteMap maps pages in a site to the links they contain to other pages on the same site. i.e. a typical "site map"
type URLTag struct {
	XMLName  xml.Name `xml:"url"`
	Location string   `xml:"loc"`
}

type URLSet struct {
	XMLName   xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
	Locations []URLTag `xml:"loc"`
}
type SiteMap map[string]bool

//const siteMapXMLStart = `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
//const siteMapXMLEnd = `</urlset>`

// GetSiteMap explores all the links that belong to the site and returns an XML site map describing the layout
func GetSiteMap(site string, depth uint) (string, error) {
	u, err := url.Parse(site)
	if err != nil {
		return "", err
	}

	// Ensure we get an absolute URL
	if !u.IsAbs() {
		return "", errors.New("unable to map site: URL is not absolute")
	}
	// Ensure we get an http/https scheme
	if !isHTTPScheme(u) {
		return "", errors.New("unable to map site: URL must be http(s) scheme")
	}

	m, err := mapSite(u, depth)
	if err != nil {
		return "", err
	}
	locations := make([]URLTag, 0, len(m))

	for k, _ := range m {
		locations = append(locations, URLTag{Location: k})
	}
	x, err := xml.MarshalIndent(URLSet{
		Locations: locations,
	}, "", "  ")
	if err != nil {
		return "", err
	}

	return string(x), nil
}

// mapSite explores all the links that belong to the site and returns a mapping of pages to links that match
func mapSite(u *url.URL, depth uint) (SiteMap, error) {
	m := SiteMap{}
	visited := map[string]bool{}
	err := doMapSite(u, depth, m, visited)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func doMapSite(u *url.URL, depth uint, m SiteMap, visited map[string]bool) error {
	// Stop if we've visited this site already
	if visited[u.String()] || depth == 0 {
		return nil
	}

	// Remember we visited this
	visited[u.String()] = true

	// Get the list of site links at this url
	r, err := getSiteLinks(u)
	if err != nil {
		return err
	}

	// Process each link
	for _, v := range r {
		// Convert to absolute URL
		nextURL, err := makeFullURL(u, v.Href)
		if err != nil {
			return err
		}

		// Only follow HTTP links to the same host
		if u.Host == nextURL.Host && isHTTPScheme(nextURL) {
			// Add to the site map
			m[nextURL.String()] = true

			// Visit this site also
			err := doMapSite(nextURL, depth-1, m, visited)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// makeFullURL parses a string URL (which could be relative) and makes it a "full" url by adding in the hostname and scheme to relative URL's
func makeFullURL(reference *url.URL, val string) (*url.URL, error) {
	u, err := url.Parse(val)
	if err != nil {
		return nil, err
	}

	// Ensure the URL has a scheme
	if u.Scheme == "" {
		u.Scheme = reference.Scheme
	}

	// Ensure the URL has a host
	if u.Host == "" {
		u.Host = reference.Host
	}

	// Don't want fragments
	u.Fragment = ""

	return u, nil
}

// getSiteLinks gets a list of links from the page to other pages on the same site
func getSiteLinks(u *url.URL) ([]links.Link, error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return links.Find(resp.Body)
}

func isHTTPScheme(u *url.URL) bool {
	return u.Scheme == "http" || u.Scheme == "https"
}
