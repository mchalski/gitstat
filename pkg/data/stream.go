package data

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

const (
	errOpen = "failed to open CsvStream on path %s"
)

// Stream is a generic data stream capable of emitting event data tuples.
type Stream interface {
	// C returns a combined data+error result channel
	// (there is no separate processing error channel)
	// on 'no more data' the channel is simply closed
	C() <-chan *StreamData
}

// StreamData is a Stream's data tuple combined with error.
type StreamData struct {
	Tuple []string
	Err   error
}

// CsvStream wraps a csv file reader.
type CsvStream struct {
	rd *csv.Reader
	c  <-chan *StreamData
}

func NewCsvStream(path string) (*CsvStream, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(errOpen, path))
	}

	c := make(chan *StreamData)
	rd := csv.NewReader(bufio.NewReader(f))

	go func() {
		defer close(c)
		for {
			rec, err := rd.Read()

			// no more records - simply close
			if rec == nil && err == io.EOF {
				break
			}

			c <- &StreamData{
				Tuple: rec,
				Err:   err,
			}
		}
	}()

	return &CsvStream{
		rd: rd,
		c:  c,
	}, nil
}

func (s *CsvStream) C() <-chan *StreamData {
	return s.c
}
