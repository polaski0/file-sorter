package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sync"
)

var wg sync.WaitGroup

func main() {
	t, o := GetPathArgs()

	dir, err := ReadDir(t)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Iterate through the array of files and copy the files from
	// the source to destination directory
	for _, f := range dir {
		if f.IsDir() {
			continue
		}

		// Create a goroutine for each file which will copy from
		// the source to destination directory concurrently
		wg.Add(1)
		go func() {
			defer wg.Done()
			file, err := os.Open(path.Join(t, f.Name()))
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			ext := filepath.Ext(f.Name())
			if ext == "" {
				ext = "misc"
			} else {
				ext = ext[1:]
			}

			dest := path.Join(o, ext)
			err = MakeDirectory(dest)
			if err != nil {
				fmt.Println(err)
				return
			}

			dfile, err := os.Create(path.Join(o, ext, f.Name()))
			if err != nil {
				fmt.Println(err)
				return
			}
			defer dfile.Close()

			rd := bufio.NewReader(file)
			wr := bufio.NewWriter(dfile)
			_, err = io.Copy(wr, rd)
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
	}

	wg.Wait()
	os.Exit(0)
}

// Read all the files in a directory
func ReadDir(path string) ([]fs.DirEntry, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	return dir, err
}

// Check if directory exists
func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// Checks if the directory exists and create a directory
// if it doesn't
func MakeDirectory(path string) error {
	ok, err := DirExists(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if ok {
		return nil
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

// Get the command-line arguments which are
// the source and destination directory respectively
func GetPathArgs() (source, destination string) {
	args := os.Args
	source, destination = "", "./tmp"

	if len(args) > 1 {
		if len(args) >= 2 {
			source = args[1]
		}

		if len(args) >= 3 {
			destination = args[2]
		}
	} else {
		fmt.Println("Target directory is required")
		os.Exit(1)
	}

	return path.Clean(source), path.Clean(destination)
}
