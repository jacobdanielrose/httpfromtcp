package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {

	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Println("read:", line)
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
