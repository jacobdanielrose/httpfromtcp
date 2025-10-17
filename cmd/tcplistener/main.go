package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
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

		linesChan := getLinesChannel(conn)

		for line := range linesChan {
			fmt.Printf("Received line: %s\n", line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)
		var currentLineContents string
		for {
			buffer := make([]byte, 8, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					lines <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			parts := strings.Split(string(buffer[:n]), "\n")
			for i := 0; i < len(parts)-1; i++ {
				currentLineContents += parts[i]
				lines <- currentLineContents
				currentLineContents = ""
			}

			currentLineContents += parts[len(parts)-1]
		}
	}()
	return lines
}
