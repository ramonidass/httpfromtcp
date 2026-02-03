package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

// `chan string` is the channel used to pass strings between different goroutines.
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
				out <- str // Send 'str' INTO channel 'out'
				str = ""
			}
			str += string(buffer)
			totalBytesRead += n

			if err == io.EOF { // End of file reached
				break
			}
			// if len(str) != 0 {
			// 	out <- str
			// }

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

	lines := getLinesReader(file)
	// 1. New variable,
	// 2. Look at the data coming out of the channel
	// 3. Figure out the type automatically, and assign it.
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}
