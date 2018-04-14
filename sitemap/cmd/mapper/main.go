package main

import (
	"fmt"
	"github.com/torbensky/gophercises/sitemap"
)

func main() {
	r, err := sitemap.GetSiteMap("https://www.calhoun.io", 2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
}
