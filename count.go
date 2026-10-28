package main

import (
	"fmt"
	"io"
	"os"
)

func countBytes(reader io.Reader) (int64, error) {
	var buffer [32 * 1024]byte
	n, err := reader.Read(buffer[:])
	return int64(n), err
}

func countTotalBytes(reader io.Reader) (int64, error) {
	var total int64
	emptyReads := 0
	
	// Allocate the 32 KiB buffer once here. 
	// Stays on the stack if it doesn't escape, or is allocated exactly once on the heap.
	var buffer [32 * 1024]byte

	for {
		// Read directly into the reusable slice window
		n, err := reader.Read(buffer[:])
		total += int64(n)
		
		if err == io.EOF {
			return total, nil
		}
		if err != nil {
			return total, err
		}

		if n == 0 {
			emptyReads++
			if emptyReads == 100 {
				return total, io.ErrNoProgress
			}
		} else {
			emptyReads = 0
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
	count, err := countTotalBytes(reader)
	if err != nil {
		return err
	}

	_, writeErr := fmt.Println(count)
	return writeErr
}

