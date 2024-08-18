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
		wg.Add(1)
		go func(src string) {
			defer wg.Done()

			fmt.Printf("Reading `%v`...\n", src)
			fds, err := s.readDirectory(src)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			for _, fd := range fds {
				fmt.Printf("Name: %v; Extension: %v; Is dotfile: %v\n",
					fd.Name(),
					path.Ext(fd.Name()),
					s.isDotFile(fd.Name()))
			}
		}(source)
	}

	wg.Wait()
	return nil
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
