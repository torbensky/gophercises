package main

import (
	"fmt"
	"log"
	"os"

	"github.com/torbensky/gophercises/links"
)

func main() {
	for _, f := range []string{"ex1.html", "ex2.html", "ex3.html", "ex4.html"} {
		fmt.Printf("links in %s:\n", f)
		printLinks(f)
	}
}

func printLinks(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	ls, err := links.Find(f)
	if err != nil {
		log.Fatal(err)
	}
	for _, l := range ls {
		fmt.Println(l)
	}
}
