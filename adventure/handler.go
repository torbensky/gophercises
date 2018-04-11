package main

import (
	"html/template"
	"log"
	"net/http"
)

type StoryPageDTO struct {
	CurrentArc string
	StoryArc
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Load initial state
	state, err := loadStateJSONFromFile(*storyPath)
	if err != nil {
		errorLoading(w)
	}

	// Convert state to our page rendering DTO
	dto := StoryPageDTO{
		state.CurrentArc,
		state.Arcs[state.CurrentArc],
	}

	// Render the HTML template
	t, err := template.ParseFiles("game.html")
	if err != nil {
		errorLoading(w)
	}
	t.Execute(w, dto)
}

func startHTTP() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func errorLoading(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Having some trouble loading stories. Check back later."))
}
