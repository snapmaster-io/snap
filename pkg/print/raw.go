package print

import "fmt"

// RawResponse prints out a raw response to the terminal
func RawResponse(response []byte) {
	fmt.Printf("Raw response:\n%s\n", string(response))
}
