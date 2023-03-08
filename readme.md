The program is meant to read/write/delete the file from a shared directory. The application handles the security of the content based on the user accessing it. Permission-based access control is set up for the application. The application takes two command line arguments - the action to perform (read, write, or delete) and the filename on which to perform the action. 

The read action reads the contents of a file in chunks and prints out each chunk up to the last newline character in that chunk. It uses a sync. Pool to manage a pool of byte slices to minimize allocations, and a sync. Mutex to synchronize access to shared variables. It uses a channel to limit the number of concurrent goroutines and reads the file in parallel by launching a new goroutine for each chunk. The read implementation will reduce the load from the in-memory loading of a file(you can read larger files now).

The write action writes user input from the standard input to a file specified by the filename parameter. If the file already exists, it prompts the user to choose whether to override or append to the file. The function uses a buffered reader to read input in smaller chunks and writes input to file in chunks of 4KB. 

The delete action deletes the file specified by the filename parameter. The minimal provided filename should be greater than three characters. 

The program also prompts the user to enter the real path to the directory where the file is located. 

Note: The program uses external packages like "bufio", "os", "sync", "io", "bytes", "github.com/google/uuid", "io/ioutil" and "path/filepath".

------------------------------------------------------------------------------------------------------------------------------------

* **BUILD THE PROJECT:** go build

* **RUNNING THE CODE:** ./note {action} {fileName}

* **How to get a real path?**
	
   **COMMAND:** pwd -P
   
   
------------------------------------------------------------------------------------------------------------------------------------