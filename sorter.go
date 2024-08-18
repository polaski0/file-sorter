package main

import (
	"fmt"
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
		go func(src string) {
			defer wg.Done()
			fds, err := s.readDirectory(src)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			for _, fd := range fds {
				if fd.IsDir() {
					continue
				}

				s.Sort(fd.Name())
			}
		}(source)
	}

	wg.Wait()
	return nil
}

func (s *Sorter) Sort(file string) {
	fname := s.getFileName(file)
	fmt.Println(fname)
}

// Generate the filename to be saved in order to avoid duplication
// which will be determined by the counter by adding a specific
// number to the current name.
func (s *Sorter) getFileName(file string) string {
	ctr := s.c.Add(file)
	ext := path.Ext(file)
	fname := file

	// If a duplicate is found, where counter is not equal to 1.
	// A count will be added to the file name.
	if ctr != 1 && !s.isDotFile(fname) {
		fname = fmt.Sprintf("%s (%d)%v",
			strings.TrimSuffix(fname, ext),
			ctr-1,
			ext)
	}

	if ctr != 1 && s.isDotFile(fname) {
		fname = fmt.Sprintf("%s (%d)",
			fname,
			ctr-1)
	}

	return fname
}

func (s *Sorter) readDirectory(path string) ([]fs.DirEntry, error) {
	fd, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

func (s *Sorter) isDotFile(path string) bool {
	return strings.HasPrefix(path, ".")
}
