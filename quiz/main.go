package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var quizfile = flag.String("csv", "problems.csv", `a csv file in the format of 'question,answer' (default "problems.csv")`)

type qresult struct {
	answer   string
	response string
}

func init() {
	flag.Parse()
}

func main() {
	fmt.Println("It's time for a quiz!")
	f, err := os.Open(*quizfile)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	// Stores quiz results
	results := map[string]qresult{}

	// Read all questions inthe quiz
	r := csv.NewReader(f)
	r.Comment = '#'
	stdr := bufio.NewReader(os.Stdin)
	var question, answer string
	for {
		// Try to read the next question
		row, err := r.Read()

		// Quiz done
		if err == io.EOF {
			quizDone(results)
			return
		}

		// Bad CSV file
		if err != nil {
			panic(err)
		}

		// Validate row
		if len(row) < 2 {
			fmt.Println("Invalid quiz format. Please ensure all rows have 2 columns. First column is question, second is answer.")
			return
		}

		// Parse question, answer and ask. Wait for answer on stdin.
		question = row[0]
		answer = row[1]
		fmt.Printf("%s ", question)
		response, _ := stdr.ReadString('\n')
		results[question] = qresult{strings.TrimSpace(answer), strings.TrimSpace(response)}
	}

}

func quizDone(results map[string]qresult) {
	total, correct := 0, 0
	for _, v := range results {
		total++
		if v.answer == v.response {
			correct++
		}
	}
	fmt.Printf("You scored %d out of %d.\n", correct, total)
}
