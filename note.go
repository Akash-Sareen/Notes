package main
import (
	"fmt"
	"os"
	"github.com/google/uuid"
	"io/ioutil"
	"strings"
	"os/user"
	"syscall"
	"strconv"

	Read "note/io/read"
	Write "note/io/write"
	Delete "note/io/delete"
	Utils "note/utils"
)

const (
	readAction = "read"
	writeAction = "write"
	deleteAction = "delete"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./program [read|write|delete] [filename]")
		os.Exit(1)
	}

	action := os.Args[1]
	filename := os.Args[2]

	var path string
	var basepath string
    fmt.Print("Enter the real path to the directory: ")
    fmt.Scanln(&path)

    hasAccess, err := checkDirectoryAccess(path)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
		os.Exit(1)
    } else if !hasAccess {
        fmt.Printf("User does not have access to directory: %v\n", path)
		os.Exit(1)
    }
	basepath = path
	path = path + "/" + filename
	switch action {
		case readAction:

			fileList, err := Utils.ListFilesWithSubstring(basepath, filename)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			if len(fileList) == 0{
				fmt.Println("No File Found")
				os.Exit(1)
			} 

			// print the file content of all files that contain the substring in the filename
			for _, file := range fileList {
				fmt.Println("File Name: " + file.Name())
				if err := Read.ReadFile(basepath + "/" + file.Name()); err != nil {
					fmt.Printf("error while reading file: %v\n", err)
				}
			}

		case writeAction:

			// generate UUID for the file name
			uuidStr := uuid.New().String()
			result := uuidStr[:4]
			fmt.Println("File Name: " + filename + "_" +result)
			fmt.Println("To save the content to the file and exit the program, press ( ctrl + d ) after typing out the content")
			path = path + "_" +result

			// Writing the string to the file
			Write.WriteToFile(path)
		
		case deleteAction:
			if len(filename) < 4 {
				fmt.Println("File Name is to short, Please enter atleat 4 character's")
				os.Exit(1)
			}

			files, err := ioutil.ReadDir(basepath)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			var isFileDeleted bool
			for _, file := range files {
				// check for end match
				if strings.HasSuffix(file.Name(), filename) {
					// deleting the file
					if err := Delete.DeleteFile(basepath + "/" + file.Name()); err != nil {
						fmt.Printf("error while deleting file: %v\n", err)
					} else {
						isFileDeleted = true
					}
				}
			}
			if(!isFileDeleted) {
				fmt.Printf("No File Found to be Deleted")
			}
		default:
			fmt.Println("Invalid action. Usage: ./program [read|write|delete] [filename]")
			os.Exit(1)
	}
}

func checkDirectoryAccess(path string) (bool, error) {
    // check if directory exists
    fileInfo, err := os.Stat(path)
    if err != nil {
        return false, fmt.Errorf("unable to access directory: %v", err)
    }

    // check if directory is readable by user
    if fileInfo.Mode().Perm()&0400 == 0 {
        return false, fmt.Errorf("no read permission for directory")
    }

    // check if directory is writable by user
    if fileInfo.Mode().Perm()&0200 == 0 {
        return false, fmt.Errorf("no write permission for directory")
    }


	// Get current user's UID
	currentUser, err := user.Current()
	if err != nil {
		return false, fmt.Errorf("Unable to access User Details")
	}
	
	// uint32(currentUser.Uid)
	uid64,err := strconv.ParseUint(currentUser.Uid, 10, 32)
	uid := uint32(uid64)

	// Get file's owner UID and group GID
	stat := fileInfo.Sys().(*syscall.Stat_t)
	ownerUid := stat.Uid
	groupGid := stat.Gid

	var permission bool

	// Check if current user is the owner or part of the group
	if uid == ownerUid {
		permission = true
	} else {
		groups, err := currentUser.GroupIds()
		if err != nil {
			return false, fmt.Errorf("Unable to access Group Details")
		}
		for _, group := range groups {
			gid, err := user.LookupGroupId(group)
			if err != nil {
				return false, fmt.Errorf("Unable to access Group Details")
			}

			gid64,err := strconv.ParseUint(gid.Gid, 10, 32)
			gid32 := uint32(gid64)
			if uint32(gid32) == groupGid {
				permission = true
			}
		}
	}
	if !permission {
		return false,fmt.Errorf("Current user is neither the owner nor part of the group that owns the folder")
	}
    return true, nil
}
