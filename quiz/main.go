package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

var quizfile = flag.String("csv", "problems.csv", `a csv file in the format of 'question,answer' (default "problems.csv")`)

type problem struct {
	question string
	answer   string
}

func init() {
	flag.Parse()
}

func main() {
	fmt.Println("It's time for a quiz!")
	f, err := os.Open(*quizfile)
	defer f.Close()
	if err != nil {
		fmt.Printf("Unable to open quiz file: %s\n", *quizfile)
		os.Exit(1)
	}

	// Stores quiz results
	results := map[problem]string{}

	// Read all questions in the quiz
	r := csv.NewReader(f)
	r.Comment = '#'
	lines, err := r.ReadAll()

	// Bad CSV file
	if err != nil {
		fmt.Println("Unable to parse CSV file")
		os.Exit(1)
	}

	// Parse the problem set
	problems := getProblems(lines)

	// Ask question for each problem
	asker := bufio.NewReader(os.Stdin)
	for i, p := range problems {
		fmt.Printf("Question %d: %s ", i+1, p.question)
		response, _ := asker.ReadString('\n')
		results[p] = strings.TrimSpace(response)
	}

	quizDone(results)
}

// getProblems gets a list of problems from CSV lines
func getProblems(lines [][]string) []problem {
	result := make([]problem, len(lines))
	for i, line := range lines {
		// Validate row
		if len(line) < 2 {
			fmt.Printf("line[%d] has an invalid problem [%s]. Please ensure all rows have at least 2 columns.\n", i, line)
			continue
		}
		result[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return result
}

// quizDone outputs the results of the quiz to the user
func quizDone(results map[problem]string) {
	correct := 0
	for p, a := range results {
		if p.answer == a {
			correct++
		}
	}
	fmt.Printf("You scored %d out of %d.\n", correct, len(results))
}
