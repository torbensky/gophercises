package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var quizFile = flag.String("csv", "problems.csv", `a csv file in the format of 'question,answer' (default "problems.csv")`)
var timeLimit = flag.Uint("limit", 30, "the quiz time limit, in seconds (default '30')")
var shuffle = flag.Bool("shuffle", false, "shuffle the problems randomly (default 'false')")

type problem struct {
	question string
	answer   string
}

func init() {
	flag.Parse()
}

func main() {
	f, err := os.Open(*quizFile)
	defer f.Close()
	if err != nil {
		fmt.Printf("Unable to open quiz file: %s\n", *quizFile)
		os.Exit(1)
	}

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
	// Stores quiz results
	results := make(map[problem]string, len(problems))

	asker := bufio.NewReader(os.Stdin)
	fmt.Println("It's time for a quiz! Hit 'enter' to start.")
	asker.ReadString('\n') // don't care about result

	responsesCh := make(chan string)
	quizTimer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	// Ask question for each problem
	for i, p := range problems {
		go func() {
			response := askQuestion(asker, i, p)
			responsesCh <- response
		}()
		select {
		case <-quizTimer.C:
			fmt.Println()
			quizDone(results, len(problems))
			return
		case response := <-responsesCh:
			results[p] = response
		}

	}

	quizDone(results, len(problems))
}

func askQuestion(r *bufio.Reader, num int, p problem) string {
	fmt.Printf("Question %d: %s ", num+1, p.question)
	response, _ := r.ReadString('\n')
	return strings.TrimSpace(response)
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

	// Shuffle the problems if the flag is set
	if *shuffle {
		for i := range result {
			j := rand.Intn(i + 1)
			result[i], result[j] = result[j], result[i]
		}
	}

	return result
}

// quizDone outputs the results of the quiz to the user
func quizDone(results map[problem]string, numProblems int) {
	correct := 0
	for p, a := range results {
		if p.answer == a {
			correct++
		}
	}
	fmt.Printf("You scored %d out of %d.\n", correct, numProblems)
}
