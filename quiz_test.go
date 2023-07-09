package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

func TestRunQuiz(t *testing.T) {
	cases := []struct {
		name          string
		pause         int
		answer        string
		expectedScore int
	}{
		{"CorrectAnswer", 0, "2", 1},
		{"IncorrectAnswer", 0, "3", 0},
		{"LateAnswer", 2, "2", 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			testRunQuiz(t, c.pause, c.answer, c.expectedScore)
		})
	}
}

func testRunQuiz(t *testing.T, pause int, answer string, expectedScore int) {
	const testFilename = "test.csv"
	var expectedOutput = fmt.Sprintf("You scored %d out of 1.", expectedScore)

	fs := fstest.MapFS{
		testFilename: {
			Data: []byte("1 + 1,2"),
		},
	}

	stdin, stdinWriter := io.Pipe()
	var stdout bytes.Buffer

	go func() {
		time.Sleep(time.Second * time.Duration(pause))
		_, _ = stdinWriter.Write([]byte(fmt.Sprintf("%s\n", answer)))
	}()

	err := runQuiz(stdin, &stdout, fs, testFilename, 1)

	if err != nil {
		t.Fatalf("Failed with error: %s", err)
	}

	result := strings.Split(stdout.String(), "\n")
	finalLine := result[len(result)-2]
	parts := strings.Split(finalLine, "=")
	output := strings.TrimSpace(parts[len(parts)-1])

	if output != expectedOutput {
		t.Fatalf("Expected the output to be '%s' but got '%s'", expectedOutput, output)
	}
}
