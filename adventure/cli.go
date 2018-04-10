package main

import (
	"errors"
	"fmt"
	"strconv"
)

func runCLIGame(s *storystate) {
	for {
		tell(*s)
		err := listOptions(*s)
		// Check if the story is over
		if err != nil {
			break
		}

		// Ask the user for a choice
		for {
			response := prompt("What would you like the little gopher to do?")

			// Check if they just want to quit
			if response == "quit" || response == "exit" {
				fmt.Println("Goodbye!")
				return
			}

			// Parse the arc number from the response
			arcNum, err := strconv.Atoi(response)
			options := s.getCurrentArc().Options
			if err == nil && arcNum > 0 && arcNum <= len(options) {
				arc := options[arcNum-1].Arc
				// Response should be a valid arc
				err = s.followArc(arc)
				if err == nil {
					break
				}
				fmt.Println(err)
			}

			fmt.Println("If you would like to quit, type 'quit' or 'exit' instead.")
		}
	}
}

// Tells the story at its current state (arc)
func tell(s storystate) {
	arc := s.getCurrentArc()
	fmt.Printf("\n\n%s\n", arc.Title)
	fmt.Println("==================================================")
	for _, v := range arc.Story {
		fmt.Printf("%s\n\n", v)
	}
}

// Asks the user which arc they'd like to follow
func prompt(msg string) string {
	fmt.Printf(msg)
	var response string
	fmt.Scanln(&response)
	return response
}

// Lists the available options to the user
// returns an error if no options available (story over)
func listOptions(s storystate) error {
	if len(s.getCurrentArc().Options) == 0 {
		return errors.New("reached end of story")
	}

	fmt.Println("What would you like to do?")
	// List the options available for the current arc
	for i, o := range s.getCurrentArc().Options {
		fmt.Printf("%d.\t%s\n", i+1, o.Text)
	}
	return nil
}
