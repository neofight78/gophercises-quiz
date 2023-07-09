package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"time"
)

func runQuiz(stdin io.Reader, stdout io.Writer, fs fs.FS, filename string, limit int) error {
	stdinReader := bufio.NewReader(stdin)

	file, err := fs.Open(filename)
	if err != nil {
		return fmt.Errorf("unable to open question file (%v): %w", filename, err)
	}

	csvReader := csv.NewReader(file)

	answers := make(chan string)

	go func() {
		for {
			answer, err := stdinReader.ReadString('\n')
			if err != nil {
				close(answers)
				break
			}
			answers <- strings.TrimSpace(answer)
		}
	}()

	timer := time.NewTimer(time.Second * time.Duration(limit))

	timedOut := false
	correct := 0
	total := 0

	for i := 0; ; i++ {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return fmt.Errorf("error reading from question file (%v): %w", filename, err)
			}
		}

		if !timedOut {
			question := record[0]
			answer := record[1]

			_, err = stdout.Write([]byte(fmt.Sprintf("Problem #%d: %s = ", i+1, question)))
			if err != nil {
				return fmt.Errorf("unable to write question to output: %w", err)
			}

			select {
			case userAnswer, ok := <-answers:
				if !ok {
					return fmt.Errorf("unable to read answer from input")
				}

				if answer == userAnswer {
					correct += 1
				}
			case <-timer.C:
				_, err = stdout.Write([]byte("\n"))
				if err != nil {
					return fmt.Errorf("unable to write to output: %w", err)
				}
				timedOut = true
			}
		}

		total += 1
	}

	_, err = stdout.Write([]byte(fmt.Sprintf("You scored %d out of %d.\n", correct, total)))
	if err != nil {
		return fmt.Errorf("unable to write score to output: %w", err)
	}

	return nil
}
