package main

import "regexp"

var re = regexp.MustCompile("(\\D*)")

func normalize(s string) string {
	return re.ReplaceAllString(s, "")
}
