package main

import (
	"errors"
	"strings"
	"testing"
)

type dataAndErrorReader struct {
	data []byte
	err  error
}

func (reader dataAndErrorReader) Read(buffer []byte) (int, error) {
	return copy(buffer, reader.data), reader.err
}

func TestCountBytes(t *testing.T) {
	count, err := countBytes(strings.NewReader("hello"))
	if err != nil {
		t.Fatalf("countBytes returned an unexpected error: %v", err)
	}
	if count != 5 {
		t.Fatalf("countBytes returned %d, want 5", count)
	}
}

func TestCountBytesReturnsPartialCountAndError(t *testing.T) {
	wantErr := errors.New("read failed")
	count, err := countBytes(dataAndErrorReader{data: []byte("abc"), err: wantErr})
	if !errors.Is(err, wantErr) {
		t.Fatalf("countBytes error = %v, want %v", err, wantErr)
	}
	if count != 3 {
		t.Fatalf("countBytes returned %d, want 3", count)
	}
}

func TestCountTotalBytesReadsUntilEOF(t *testing.T) {
	input := strings.Repeat("a", 100_000)
	count, err := countTotalBytes(strings.NewReader(input))
	if err != nil {
		t.Fatalf("countTotalBytes returned an unexpected error: %v", err)
	}
	if count != int64(len(input)) {
		t.Fatalf("countTotalBytes returned %d, want %d", count, len(input))
	}
}

func TestCountTotalBytesReturnsPartialCountAndError(t *testing.T) {
	wantErr := errors.New("read failed")
	count, err := countTotalBytes(dataAndErrorReader{data: []byte("abc"), err: wantErr})
	if !errors.Is(err, wantErr) {
		t.Fatalf("countTotalBytes error = %v, want %v", err, wantErr)
	}
	if count != 3 {
		t.Fatalf("countTotalBytes returned %d, want 3", count)
	}
}
