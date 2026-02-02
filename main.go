package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

// <-chan string is the channel used to pass strings
// between different goroutines.
// <- Arrow is on the left (Receive-Only operator)
func getLinesReader(f io.ReadCloser) <-chan string { // (ReadCloser is an interface that accept Read() and Close() methods)
	out := make(chan string, 1)
	go func() {
		defer f.Close()  // Ensure the file is closed when finished
		defer close(out) // Ensure the channel is closed when finished
		totalBytesRead := 0
		str := ""
		fmt.Println("Content of messages.txt")
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				break
			}
			buffer = buffer[:n]
			if i := bytes.IndexByte(buffer, '\n'); i != -1 {
				str += string(buffer[:i])
				buffer = buffer[i+1:]
				fmt.Printf("read: %s\n", str)
				str = ""
			}
			str += string(buffer)
			totalBytesRead += n

			if err == io.EOF { // End of file reached
				break
			}
			// fmt.Printf("\nTotal bytes read: %d\n", totalBytesRead)
		}
	}()
	return out
}

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("Error closing file: %v", closeErr)
		}
	}() // Invoke the anonymous function

	lines := getLinesReader(f)
}
