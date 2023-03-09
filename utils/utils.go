package utils

import (
	"os"
	"io/ioutil"
	"strings"
	// "path/filepath"
)

func ListFilesWithSubstring(dirPath, substr string) ([]os.FileInfo, error) {
	// read all files in the specified directory and its subdirectories
	fileList, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	// filter the file list to include only files that contain the substring
	var filteredList []os.FileInfo
	for _, file := range fileList {
		if file.Mode().IsRegular() && strings.Contains(file.Name(), substr) {
			filteredList = append(filteredList, file)
		}
	}

	// Assuming in future the code can used to save notes according to category directory.
	// recursively search subdirectories for files that contain the substring
	// subdirs, err := ioutil.ReadDir(dirPath)
	// if err != nil {
	// 	return nil, err
	// }
	// for _, subdir := range subdirs {
	// 	if subdir.IsDir() {
	// 		sublist, err := ListFilesWithSubstring(filepath.Join(dirPath, subdir.Name()), substr)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		filteredList = append(filteredList, sublist...)
	// 	}
	// }

	return filteredList, nil
}