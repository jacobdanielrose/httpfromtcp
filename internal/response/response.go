package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacobdanielrose/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

var statusMessage = map[StatusCode]string{
	StatusOK:                  "OK",
	StatusBadRequest:          "Bad Request",
	StatusInternalServerError: "Internal Server Error",
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLineStr := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusMessage[statusCode])
	_, err := w.Write([]byte(statusLineStr))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		_, err := w.Write(fmt.Appendf([]byte{}, "%s: %s\r\n", key, val))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}
