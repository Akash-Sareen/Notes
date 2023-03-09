package write

import (
	"bufio"
	"fmt"
	"os"
	"errors"
	"io"
)

const (
	bufferSize = 1 * 1024 * 1024
)


// writeToFile is a function that writes user input from the standard input to a file specified by the filename parameter.
// If the file already exists, it prompts the user to choose whether to override or append to the file.
// The function uses a buffered reader to read input in smaller chunks, and writes input to file in chunks of 1Mb.
// It returns nothing, but prints a success message when the input is saved to the file successfully.
func WriteToFile(filename string) {
	var mode int

	if _, err := os.Stat(filename); err == nil {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("File '%s' already exists, do you want to override it? (y/n): ", filename)
		input, _ := reader.ReadString('\n')
		if input == "y\n" {
			mode = os.O_TRUNC
		} else {
			mode = os.O_APPEND
		}
	} else {
		mode = os.O_CREATE
	}

	file, err := os.OpenFile(filename, mode|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a buffered reader to read input in smaller chunks
	reader := bufio.NewReader(os.Stdin)

	// Write input to file in chunks
	buffer := make([]byte, bufferSize)
	for {
		bytesRead, err := reader.Read(buffer)
		if err != nil && !errors.Is(err, io.EOF) {
			fmt.Println("Error reading input:", err)
			return
		}
		if bytesRead == 0 {
			break
		}
		if _, err := file.Write(buffer[:bytesRead]); err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	fmt.Printf("Input saved to '%s' successfully!\n", filename)
}
