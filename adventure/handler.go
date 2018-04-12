package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
)

type StoryHandler struct {
	story StoryState
}

type StoryPageDTO struct {
	CurrentArc string
	StoryArc
}

var gotoPathRegex = regexp.MustCompile(`^/([^"]+)/goto/([^"]+)$`)

func (s StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Let's do a really basic implementation only using the core http library, no 3rd-party routers :O
	path := r.URL.Path
	if path == "/" {
		s.start(w, r)
	} else if gotoPathRegex.MatchString(path) {
		s.goToArc(w, r)
	} else {
		// TODO: This could be nicer. I think we should communicate this is a 404, but also be more helpful.
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}
}

func (s StoryHandler) start(w http.ResponseWriter, r *http.Request) {
	renderStory(w, s.story)
}

func (s StoryHandler) goToArc(w http.ResponseWriter, r *http.Request) {
	// This should match, validation should have happened in the router
	matches := gotoPathRegex.FindStringSubmatch(r.URL.Path)

	// Safety first!
	if matches == nil || len(matches) != 3 {
		log.Printf("goToArc received an invalid path: %s\n", r.URL.Path)
		errorLoading(w)
		return
	}

	fromArc, toArc := matches[1], matches[2]
	story := s.story // copy for thead-safety. we gonna mutate story

	// Set current (from) arc
	err := story.setArc(fromArc)
	if err != nil {
		errorBadRequest(w, err)
		return
	}

	err = story.followArc(toArc)
	if err != nil {
		errorBadRequest(w, err)
		return
	}

	renderStory(w, story)
}

func startHTTP() {
	// Load initial state
	state, err := loadStateJSONFromFile(*storyPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":8080", &StoryHandler{
		story: *state,
	}))
}

func errorBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

func errorLoading(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Having some trouble loading stories. Check back later."))
}

func renderStory(w http.ResponseWriter, story StoryState) {
	// Convert state to our page rendering DTO
	dto := StoryPageDTO{
		story.CurrentArc,
		story.Arcs[story.CurrentArc],
	}

	// Render the HTML template
	t, err := template.ParseFiles("game.html")
	if err != nil {
		errorLoading(w)
	}
	t.Execute(w, dto)

}
