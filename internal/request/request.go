package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf, _ := io.ReadAll(reader)

	requestLine, err := parseRequestLine(buf)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: requestLine,
	}, nil
}

func parseRequestLine(buf []byte) (RequestLine, error) {
	buf_string := string(buf)
	requestLinestring := strings.Split(buf_string, "\r\n")[0]
	splitRequestLine := strings.Split(requestLinestring, " ")

	if len(splitRequestLine) != 3 {
		fmt.Print(len(splitRequestLine))
		return RequestLine{}, errors.New("Invalid request line")
	}

	method := splitRequestLine[0]

	if !isValidMethod(method) {
		return RequestLine{}, errors.New("Not a valid method")
	}

	requestTarget := splitRequestLine[1]

	httpVersion := strings.Split(splitRequestLine[2], "/")[1]
	if !isValidHttpVersion(httpVersion) {
		return RequestLine{}, errors.New("Not a valid HTTP version")
	}

	return RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: requestTarget,
		Method:        method,
	}, nil
}

func isValidMethod(method string) bool {
	for _, letter := range method {
		if !unicode.IsLetter(letter) || !unicode.IsUpper(letter) {
			return false
		}
	}
	return true
}

func isValidHttpVersion(versionString string) bool {
	if versionString != "1.1" {
		return false
	}
	return true
}
