package delete

import (
	"os"
	"fmt"
)

func DeleteFile(filePath string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", filePath)
	}

	// Remove the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("unable to delete file %s: %v", filePath, err)
	}

	return nil
}