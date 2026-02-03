package main

import (
	"bufio"
	"io"
	"log"
)

func getLinesScanner(f io.ReadCloser) <-chan string {
	out := make(chan string) // Unbuffered is fine here, or small buffer
	go func() {
		defer f.Close()
		defer close(out)

		// Scanner splits by lines by default
		scanner := bufio.NewScanner(f)

		// Scan() returns true as long as there is data to read
		for scanner.Scan() {
			// scanner.Text() extracts the string from the current line
			out <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Error reading file: %v", err)
		}
	}()
	return out
}
