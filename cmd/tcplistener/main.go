package main

import (
	"fmt"
	"net"

	"github.com/jacobdanielrose/httpfromtcp/internal/request"
)

const port = ":42069"

func main() {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("listening for TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			continue
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		request, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			continue
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range request.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
		fmt.Println("Body:")
		fmt.Printf("%s\n", request.Body)

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}

}
