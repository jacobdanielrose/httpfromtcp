package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/jacobdanielrose/httpfromtcp/internal/request"
	"github.com/jacobdanielrose/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(port int, h Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
	}

	go server.listen(h)
	return server, nil

}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen(h Handler) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn, h)
	}
}

func (s *Server) handle(conn net.Conn, h Handler) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("There was an error parsing request: %v\n", err)
		headers := response.GetDefaultHeaders(0)
		if err := response.WriteHeaders(conn, headers); err != nil {
			fmt.Printf("error: %v\n", err)
		}
		response.WriteStatusLine(conn, response.StatusInternalServerError)
	}

	bodyBuf := bytes.Buffer{}
	hErr := h(&bodyBuf, req)

	if hErr != nil {
		msg := []byte(hErr.Message)
		response.WriteStatusLine(conn, hErr.StatusCode)
		headers := response.GetDefaultHeaders(len(msg))
		if err := response.WriteHeaders(conn, headers); err != nil {
			fmt.Printf("error: %v\n", err)
		}
		if _, err := conn.Write(msg); err != nil {
			fmt.Printf("error: %v\n", err)
		}
		return
	}
	response.WriteStatusLine(conn, response.StatusOK)
	headers := response.GetDefaultHeaders(len(bodyBuf.Bytes()))
	if err := response.WriteHeaders(conn, headers); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	bodyBuf.WriteTo(conn)
	if err != nil {
		fmt.Printf("error : %v\n", err)
		return
	}
}
