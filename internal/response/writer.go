package response

import (
	"fmt"
	"io"

	"github.com/jacobdanielrose/httpfromtcp/internal/headers"
)

type writerState int

const (
	writingStatus writerState = iota
	writingHeaders
	writingBody
	writingTrailers
)

type Writer struct {
	writer io.Writer
	state  writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		state:  writingStatus,
		writer: w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writingStatus {
		return fmt.Errorf("cannot write status line in state %d", w.state)
	}
	defer func() { w.state = writingHeaders }()
	_, err := w.writer.Write(getStatusLine(statusCode))
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != writingHeaders {
		return fmt.Errorf("cannot write headers in state %d", w.state)
	}
	defer func() { w.state = writingBody }()
	for key, val := range headers {
		_, err := w.writer.Write(fmt.Appendf([]byte{}, "%s: %s\r\n", key, val))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteTrailers(trailers headers.Headers) error {
	if w.state != writingTrailers {
		return fmt.Errorf("cannot write trailers in state %d", w.state)
	}
	defer func() { w.state = writingBody }()
	for key, val := range trailers {
		_, err := w.writer.Write(fmt.Appendf([]byte{}, "%s: %s\r\n", key, val))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writingBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.state)
	}
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != writingBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.state)
	}
	chunkSize := len(p)

	nTotal := 0
	n, err := fmt.Fprintf(w.writer, "%x\r\n", chunkSize)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write(p)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return nTotal, err
	}
	nTotal += n
	return nTotal, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != writingBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.state)
	}
	n, err := w.writer.Write([]byte("0\r\n"))
	if err != nil {
		return n, err
	}
	w.state = writingTrailers
	return n, nil
}
