package response

import "fmt"

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

var StatusMessage = map[StatusCode]string{
	StatusOK:                  "OK",
	StatusBadRequest:          "Bad Request",
	StatusInternalServerError: "Internal Server Error",
}

func getStatusLine(statusCode StatusCode) []byte {
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, StatusMessage[statusCode]))
}
