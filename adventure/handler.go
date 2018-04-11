package main

import (
	"fmt"
	"log"
	"net/http"
)

type StoryPageDTO struct {
	CurrentArc string
	StoryArc
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func startHTTP() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
