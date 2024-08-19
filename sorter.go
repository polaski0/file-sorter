package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"sync"
)

type Sorter struct {
	c    *Counter
	src  []string
	dest string
}

var wg sync.WaitGroup

func NewSorter(src []string, dest string) *Sorter {
	c := NewCounter()
	return &Sorter{
		c:    c,
		src:  src,
		dest: dest,
	}
}

func (s *Sorter) Start() error {
	for _, source := range s.src {
		fmt.Printf("Reading `%v`...\n", source)
		wg.Add(1)
		go func() {
			defer wg.Done()
			fds, err := readDirectory(source)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			for _, fd := range fds {
				if fd.IsDir() {
					continue
				}

				fmt.Printf("Sorting %v from %v...\n", fd.Name(), source)
				err := s.Sort(fd.Name(), source)
				if err != nil {
					fmt.Printf("Error %v", err)
					return
				}
				fmt.Printf("Sorted %v\n", fd.Name())
			}
		}()
	}

	wg.Wait()
	return nil
}

func (s *Sorter) Sort(file, source string) error {
	var dest string
	ext := path.Ext(file)
	done := make(chan bool)
	error := make(chan error)

	if isDotFile(file) || ext == "" {
		dest = path.Join(s.dest, "misc")
	} else {
		dest = path.Join(s.dest, ext[1:])
	}

	err := os.MkdirAll(dest, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	go s.createFile(file, source, dest, error, done)

	for {
		select {
		case err := <-error:
			return err
		case <-done:
			return nil
		}
	}
}

// the fuck?
func (s *Sorter) createFile(file string, source string, dest string, error chan error, done chan bool) {
	fname := s.getFileName(file) // what
	exist, err := isFileExists(fname, dest)
	if err != nil && !os.IsNotExist(err) {
		error <- err
	}

	if exist {
		go s.createFile(file, source, dest, error, done)
		return
	}

	f, err := os.Open(path.Join(source, file))
	if err != nil {
		error <- err
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	wr, err := os.Create(path.Join(dest, fname))
	if err != nil {
		error <- err
	}
	defer wr.Close()

	_, err = io.Copy(wr, rd)
	if err != nil {
		error <- err
	}

	done <- true
}

func isFileExists(file, dir string) (bool, error) {
	_, err := os.Stat(path.Join(dir, file))
	if err != nil {
		return false, err
	}
	return true, nil
}

// Generate the filename to be saved in order to avoid duplication
// which will be determined by the counter by adding a specific
// number to the current name.
func (s *Sorter) getFileName(file string) string {
	ctr := s.c.Add(file)
	ext := path.Ext(file)

	// If a duplicate is found, where counter is not equal to 1.
	// A count will be added to the file name.
	if ctr != 1 && !isDotFile(file) {
		file = fmt.Sprintf("%s (%d)%v",
			strings.TrimSuffix(file, ext),
			ctr-1,
			ext)
	}

	if ctr != 1 && isDotFile(file) {
		file = fmt.Sprintf("%s (%d)",
			file,
			ctr-1)
	}

	return file
}

func readDirectory(path string) ([]fs.DirEntry, error) {
	fd, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

func isDotFile(path string) bool {
	return strings.HasPrefix(path, ".")
}
