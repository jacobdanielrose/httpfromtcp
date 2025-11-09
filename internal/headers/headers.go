package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		// the empty line
		// headers are done, consume the CRLF
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	if len(parts) != 2 {
		return 0, false, fmt.Errorf("malformed header: missing colon")
	}

	key := strings.ToLower(string(parts[0]))

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)

	matched, err := regexp.MatchString("^[A-Za-z0-9_!#$%&`'+._|~-]+$", key)
	if err != nil {
		return 0, false, err
	}
	if !matched {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	existingVal, exists := h[key]
	if exists {
		h.Set(key, existingVal+", "+string(value))
	} else {
		h.Set(key, string(value))
	}

	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{
			v,
			value,
		}, ", ")
	}
	h[key] = value
}

func (h Headers) Override(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	v, ok := h[key]
	return v, ok
}

func (h Headers) Delete(key string) (string, bool) {
	key = strings.ToLower(key)
	v, ok := h[key]
	delete(h, key)
	return v, ok
}
