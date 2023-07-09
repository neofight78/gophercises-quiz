package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	csv := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	limit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	err := runQuiz(os.Stdin, os.Stdout, os.DirFS("."), *csv, *limit)
	if err != nil {
		log.Fatalf("an unhandled error occured: %v", err)
	}
}
