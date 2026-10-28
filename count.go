package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
	var filename string
	var countBytes, countLines, countWords, countChars, hasFlags bool

	// =========================================================================
	// PHASE 1: Parse arguments sequentially
	// =========================================================================
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-c":
			countBytes = true
			hasFlags = true
		case "-l":
			countLines = true
			hasFlags = true
		case "-w":
			countWords = true
			hasFlags = true
		case "-m":
			countChars = true
			hasFlags = true
		default:
			// If it starts with a hyphen, it's an unrecognized option
			if len(arg) > 0 && arg[0] == '-' {
				return fmt.Errorf("ccwc: invalid option -- '%s'", arg[1:])
			}
			// If filename is already populated, the user passed multiple files
			if filename != "" {
				return fmt.Errorf("ccwc: extra operand '%s' (only one file at a time is supported)", arg)
			}
			filename = arg
		}
	}

	// =========================================================================
	// PHASE 2: Apply default behaviors if no flags were provided
	// =========================================================================
	if !hasFlags {
		countLines = true
		countWords = true
		countBytes = true
	}

	// =========================================================================
	// PHASE 3: Resolve input source stream (File vs. Stdin)
	// =========================================================================
	var reader io.Reader
	if filename != "" {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	// =========================================================================
	// PHASE 4: Compute metrics in a single stream pass
	// =========================================================================
	metrics, err := countTotalMetrics(reader)
	if err != nil {
		return err
	}

	// =========================================================================
	// PHASE 5: Print matching results in standard Unix wc layout order
	// =========================================================================
	if countLines {
		fmt.Printf(" %7d", metrics.Lines)
	}
	if countWords {
		fmt.Printf(" %7d", metrics.Words)
	}
	if countBytes {
		fmt.Printf(" %7d", metrics.Bytes)
	}
	if countChars {
		fmt.Printf(" %7d", metrics.Chars)
	}
	if filename != "" {
		fmt.Printf(" %s", filename)
	}
	fmt.Println()

	return nil
}
