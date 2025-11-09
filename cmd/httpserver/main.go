package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jacobdanielrose/httpfromtcp/internal/request"
	"github.com/jacobdanielrose/httpfromtcp/internal/response"
	"github.com/jacobdanielrose/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		proxyHandler(w, req)
		return
	}
	handler200(w, req)
}

func proxyHandler(w *response.Writer, req *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	headers := response.GetDefaultHeaders(0)
	headers.Override("Content-Type", "text/html")
	headers.Delete("Content-Length")
	headers.Set("Transfer-Encoding", "chunked")
	w.WriteHeaders(headers)

	resp, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", path))
	if err != nil {
		log.Printf("Could not proxy request: %v", err)
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 32)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("Could not read response body: %v", err)
			return
		}
		if n == 0 {
			break
		}
		// buf := fmt.Appendf(buf[:n], "n: %d", n)
		w.WriteChunkedBody(buf[:n])
	}
	w.WriteChunkedBodyDone()
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusBadRequest)
	body := returnHTML(
		response.StatusBadRequest,
		response.StatusMessage[response.StatusBadRequest],
		"Your request honestly kinda sucked.",
	)
	headers := response.GetDefaultHeaders(len(body))
	headers.Override("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(body)
}
func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusInternalServerError)
	body := returnHTML(
		response.StatusInternalServerError,
		response.StatusMessage[response.StatusInternalServerError],
		"Okay, you know what? This one is on me.",
	)
	headers := response.GetDefaultHeaders(len(body))
	headers.Override("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	body := returnHTML(
		response.StatusOK,
		"Success!",
		"Your request was an absolute banger.",
	)
	headers := response.GetDefaultHeaders(len(body))
	headers.Override("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func returnHTML(statusCode response.StatusCode, statusMsg string, messageBody string) []byte {
	statusStr := response.StatusMessage[statusCode]
	htmlStr := fmt.Sprintf(`<html>
<head>
<title>%d %s</title>
</head>
<body>
<h1>%s</h1>
<p>%s</p>
</body>
</html>`,
		statusCode,
		statusStr,
		statusMsg,
		messageBody,
	)
	return []byte(htmlStr)
}
