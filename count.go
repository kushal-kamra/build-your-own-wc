package main

import (
	"fmt"
	"io"
	"os"
	"bufio"
	"unicode"

)

// Metrics stores the aggregated counts computed during a single pass over a stream.
type Metrics struct {
	Lines int64
	Words int64
	Bytes int64
	Chars int64
}

// countTotalMetrics processes the provided reader exactly once using an efficient
// buffered stream scanner. It dynamically tracks lines, words, bytes, and characters
// simultaneously without risking multi-byte encoding corruption or unnecessary allocations.
func countTotalMetrics(reader io.Reader) (Metrics, error) {
	// Initialize struct fields to zero automatically
	var metrics Metrics
	
	// Wrap the reader to protect against multi-byte slice splitting on boundaries
	bufferedReader := bufio.NewReader(reader)
	
	// Finite state machine tracker for word boundaries
	inWord := false

	for {
		// Decode one full Unicode code point (rune) safely from the stream buffer
		r, size, err := bufferedReader.ReadRune()
		if err == io.EOF {
			// Stream ended successfully; return accumulated metrics with a nil error
			return metrics, nil
		}
		if err != nil {
			// A genuine system or network failure occurred; return what we have so far
			return metrics, err
		}

		// 1. Increment characters and exact consumed byte sizes
		metrics.Chars++
		metrics.Bytes += int64(size)

		// 2. Count newlines
		if r == '\n' {
			metrics.Lines++
		}

		// 3. Process the word-counting state engine using comprehensive Unicode tracking
		if unicode.IsSpace(r) {
			inWord = false
		} else {
			if !inWord {
				metrics.Words++
				inWord = true
			}
		}
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	// 1. Declare a single interface variable to hold our input source
	var reader io.Reader

	// 2. Inspect the arguments to determine the source
	if len(os.Args) > 1 {
		// Filename supplied: open the file
		filename := os.Args[1]
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		// Ensure the file is closed when run() exits
		defer file.Close()
		
		// Assign the file to our unified reader
		reader = file
	} else {
		// No filename supplied: fallback to standard input
		reader = os.Stdin
	}

	// 3. Execution pipeline (Unified Logic)
	// Both paths converge here seamlessly without code duplication
	metrics, err := countTotalMetrics(reader)
	if err != nil {
		return err
	}

	_, writeErr := fmt.Println(metrics)
	return writeErr
}

