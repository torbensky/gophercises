package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type storystate struct {
	currentArc string
	arcs       map[string]StoryArc
}

type StoryArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

var storyPath = flag.String("json", "gopher.json", "a path to a valid story JSON file (default 'story.json')")

func init() {
	flag.Parse()
}

func main() {
	state, err := loadStateJson(*storyPath)
	if err != nil {
		fmt.Printf("Unable to load story at: %s\n", *storyPath)
		os.Exit(1)
	}

	runCLIGame(state)
}

// 1. Want to load the JSON file
//	- consider doing this in a modular way so other format could be supported
// 2. Want to parse the JSON contents into our CYOA structure, probably using structs
// 3. Want to start a game using the initial state parsed from the JSON file
//	- 'Intro' is the initial state
//	- Each state will provide a list of transitions

func loadStateJson(path string) (*storystate, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var arcs map[string]StoryArc
	err = json.Unmarshal(b, &arcs)
	if err != nil {
		return nil, err
	}

	return &storystate{
		currentArc: "intro",
		arcs:       arcs,
	}, nil
}

func runCLIGame(s *storystate) {
	for {
		tell(*s)
		err := listOptions(*s)
		// Check if the story is over
		if err != nil {
			break
		}

		// Continue to nag prompt until a valid arc is received
		for {
			arc := promptForArc()
			err = s.followArc(arc)
			if err == nil {
				break
			}
		}
	}
}

func tell(s storystate) {
	arc := s.getCurrentArc()
	fmt.Printf("\n\n%s\n\n", arc.Title)
	for _, v := range arc.Story {
		fmt.Printf("%s\n\n", v)
	}
}

func promptForArc() string {
	fmt.Printf("\nChoose a valid option from the names above: ")
	var response string
	fmt.Scanln(&response)
	return response
}

func listOptions(s storystate) error {
	if len(s.getCurrentArc().Options) == 0 {
		return errors.New("reached end of story")
	}

	fmt.Println("What would you like to do?")
	// List the options available for the current arc
	for _, o := range s.getCurrentArc().Options {
		fmt.Printf("\n%s:  %s\n", o.Arc, o.Text)
	}
	return nil
}

func (s *storystate) setArc(arc string) error {
	arcValid := false
	for name, _ := range s.arcs {
		if name == arc {
			arcValid = true
			break
		}
	}

	if !arcValid {
		return errors.New("invalid story arc: '" + arc + "'")
	}

	s.currentArc = arc
	return nil
}

func (s *storystate) followArc(arc string) error {
	arcValid := false
	for _, s := range s.getCurrentArc().Options {
		if s.Arc == arc {
			arcValid = true
			break
		}
	}
	if !arcValid {
		return errors.New("invalid story arc: '" + arc + "'")
	}

	s.currentArc = arc

	return nil
}

func (s *storystate) getCurrentArc() StoryArc {
	return s.arcs[s.currentArc]
}
