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

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writingBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.state)

	}
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	chunkLen := fmt.Sprintf("%x\r\n", len(p))
	n1, err := w.writer.Write([]byte(chunkLen))
	n2, err := w.writer.Write(p)
	n3, err := w.writer.Write([]byte("\r\n"))
	return n1 + n2 + n3, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.writer.Write([]byte("0\r\n\r\n"))
}
