package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type StoryState struct {
	CurrentArc string
	Arcs       map[string]StoryArc
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

func loadStateJSONFromFile(path string) (*StoryState, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return loadStateJSON(b)
}

func loadStateJSON(b []byte) (*StoryState, error) {
	var arcs map[string]StoryArc
	err := json.Unmarshal(b, &arcs)
	if err != nil {
		return nil, err
	}

	return &StoryState{
		CurrentArc: initialStoryArc,
		Arcs:       arcs,
	}, nil
}

func (s *StoryState) setArc(arc string) error {
	if _, ok := s.Arcs[arc]; ok {
		s.CurrentArc = arc
		return nil
	}

	return errors.New("invalid story arc: '" + arc + "'")
}

func (s *StoryState) followArc(arc string) error {
	for _, sa := range s.getCurrentArc().Options {
		if sa.Arc == arc {
			s.CurrentArc = arc
			return nil
		}
	}
	return errors.New("invalid story arc: '" + arc + "'")
}

func (s *StoryState) getCurrentArc() StoryArc {
	return s.Arcs[s.CurrentArc]
}
