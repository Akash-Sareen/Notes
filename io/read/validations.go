package read

import (
	"fmt"
	"os"
)

// Validation for the filename, used in the read text from file function
func validateFilename(filename string) error {
	if filename == "" {
		return fmt.Errorf("empty filename")
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filename)
	}
	return nil
}