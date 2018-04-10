package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type storystate struct {
	currentArc string
	arcs       map[string]StoryArc
}

// StoryArc is just that - the arc of a story. Basically, a chapter.
type StoryArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

const initialStoryArc = "intro"

func loadStateJSON(path string) (*storystate, error) {
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
		currentArc: initialStoryArc,
		arcs:       arcs,
	}, nil
}

func (s *storystate) setArc(arc string) error {
	arcValid := false
	for name := range s.arcs {
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
