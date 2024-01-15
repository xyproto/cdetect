package ainur

import (
	"errors"
	"io"
)

// StreamReader is intended to be used to search in a streaming manner.
type StreamReader struct {
	buf []byte
	r   io.Reader
}

// NewStreamReader creates a new StreamReader.
func NewStreamReader(r io.Reader, bufferSize int) (*StreamReader, error) {
	if bufferSize%2 != 0 {
		return nil, errors.New("buffer size must be even")
	}

	return &StreamReader{
		r:   r,
		buf: make([]byte, bufferSize),
	}, nil
}

// Next reads the next half byte buffer from the stream.
// The reason for reading half is so that we can build something like:
//
// tail my-file | grep "something"
//
// Where "something" may be between two reads.
func (r *StreamReader) Next() ([]byte, error) {
	copy(r.buf, r.buf[len(r.buf)/2:])

	half := len(r.buf) / 2
	n, err := io.ReadFull(r.r, r.buf[half:])
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			return r.buf[:half+n], nil
		}
		return nil, err
	}

	return r.buf[:half+n], nil
}
