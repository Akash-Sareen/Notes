package read

import (
	"fmt"
	"sync"
	"io"
	"bytes"
	"os"
)

const (
	chunkSize   = 10 * 1024 * 1024 // 10Mb size of chunk
)

// readFile reads the contents of a file in chunks and prints out each chunk up to the last newline character in that chunk.
// It uses a sync.Pool to manage a pool of byte slices to minimize allocations, and a sync.Mutex to synchronize access to shared variables.
// It uses a channel to limit the number of concurrent goroutines, and reads the file in parallel by launching a new goroutine for each chunk. It might be an over do for this use case, but the same code can be used for threaded application.
// The method takes a filename as a string and returns an error if the file cannot be opened or read.
func ReadFile(filename string) error {
	if err := validateFilename(filename); err != nil {
		return err
	}

	// Open the file for reading
	file, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Get the file info to determine its size
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	// Determine the size of the file
	fileSize := info.Size()

	// Create a pool of buffers to read the file in chunks
	var pool = sync.Pool{
		New: func() interface{} {
			return make([]byte, chunkSize)
		},
	}

	// Create a semaphore to limit the number of concurrent reads
	var sem = make(chan struct{}, 100)

	// Create a mutex to protect access to the file offset and current chunk
	var mutex sync.Mutex
	var offset int64
	var chunk []byte
	offset = 0

	// Read the file in chunks using the pool of buffers and the semaphore
	for offset < fileSize {
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()

			mutex.Lock()
			defer mutex.Unlock()

			chunk, err := getChunk(file, offset, pool)
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}

			if chunk == nil {
				return
			}

			offset += int64(len(chunk))

			lastNewline := bytes.LastIndexByte(chunk, '\n')
			if lastNewline == -1 {
				fmt.Printf("%s", chunk)
			} else 	if (offset - int64(len(chunk))) != fileSize {
				fmt.Printf("%s", chunk[:fileSize])
			} else {
				fmt.Printf("%s", chunk[:lastNewline+1])
			}
		}()
	}
	
	// Print any remaining chunk
	if len(chunk) > 0 {
		fmt.Printf("%s\n", chunk)
	}

	// Wait for all reads to complete before returning
	waitForReads(sem)

	file.Close()
	return nil
}

func getChunk(file *os.File, offset int64, pool sync.Pool) ([]byte, error) {
	buf := pool.Get().([]byte)
	defer pool.Put(buf)

	chunk := buf[:chunkSize]
	n, err := file.ReadAt(chunk, offset)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("error while reading file at offset %d: %v", offset, err)
	}

	if n == 0 {
		return nil, nil
	}

	return chunk[:n], nil
}

func waitForReads(sem chan struct{}) {
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
}
