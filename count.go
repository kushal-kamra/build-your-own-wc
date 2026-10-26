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
	for {
		count, err := countBytes(reader)
		total += count
		if err == io.EOF {
			return total, nil
		}
		if err != nil {
			return total, err
		}
		if count == 0 {
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
	file, err := os.Open("test.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	count, err := countTotalBytes(file)
	if _, writeErr := fmt.Println(count); writeErr != nil {
		return writeErr
	}
	return err
}
