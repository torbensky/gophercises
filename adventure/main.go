package main

import (
	"flag"
	"fmt"
	"os"
)

var storyPath = flag.String("json", "gopher.json", "a path to a valid story JSON file (default 'story.json')")

func init() {
	flag.Parse()
}

func main() {
	startHTTP()
	state, err := loadStateJSON(*storyPath)
	if err != nil {
		fmt.Printf("Unable to load story at: %s\n", *storyPath)
		os.Exit(1)
	}

	runCLIGame(state)
}
