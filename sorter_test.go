package main

import (
	"fmt"
	"os"
	"path"
	"sync"
	"testing"
)

var twg sync.WaitGroup

const dir = "./tmp"

// Create test files on the current directory
func setup() {
	err := os.Mkdir(dir, 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error: %v", err)
		return
	}

	type File struct {
		filename string
		value    string
	}

	files := map[string][]File{
		"dir1": []File{
			{
				filename: "test.txt",
				value:    "The quick brown fox jumps over the lazy dog",
			},
			{
				filename: "test (1).txt",
				value:    "Another quick brown fox jumps over the lazy dog",
			},
			{
				filename: "test (2).txt",
				value:    "iwashere",
			},
			{
				filename: "README.md",
				value:    "# Hello\nworld",
			},
			{
				filename: "Fizz.md",
				value:    "# Buzz",
			},
		},
		"dir2": []File{
			{
				filename: "test.txt",
				value:    "The quick brown fox jumps over the lazy dog",
			},
			{
				filename: "test (1).txt",
				value:    "Another quick brown fox jumps over the lazy dog",
			},
			{
				filename: "README.md",
				value:    "# Goodbye\nworld",
			},
		},
	}

	for k, v := range files {
		dest := path.Join(dir, k)
		if err := os.Mkdir(dest, 0755); err != nil {
			continue
		}

		for _, f := range v {
			twg.Add(1)
			go func() {
				defer twg.Done()
				err := os.WriteFile(path.Join(dest, f.filename), []byte(f.value), 0644)
				if err != nil {
					fmt.Printf("Error writing file %v, %v", f.filename, err)
					return
				}
			}()
		}
	}

	twg.Wait()
}

// Remove test files
func teardown() {
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
}

func TestSort(t *testing.T) {
	setup()
	defer teardown()

	dest := path.Join(dir, "out")
	src := []string{
		path.Join(dir, "dir1"),
		path.Join(dir, "dir2"),
	}

	sr := NewSorter(src, dest)
	err := sr.Start()
	if err != nil {
		t.Errorf("Error %v\n", err)
	}
}

func TestFileExists(t *testing.T) {
	setup()
	defer teardown()

	testCases := []struct {
		dir      string
		file     string
		expected bool
	}{
		{
			dir:      path.Join(dir, "dir1"),
			file:     "test.txt",
			expected: true,
		},
		{
			dir:      path.Join(dir, "dir1"),
			file:     "test (1).txt",
			expected: true,
		},
		{
			dir:      path.Join(dir, "dir1"),
			file:     "invalid-file.md",
			expected: false,
		},
	}

	for i := range testCases {
		exists, err := isFileExists(testCases[i].file, testCases[i].dir)
		if err != nil && !os.IsNotExist(err) {
			t.Errorf("Error %v", err)
		}

		if exists != testCases[i].expected {
			t.Errorf("File: %v; Found %v, expected %v\n",
				testCases[i].file,
				exists,
				testCases[i].expected)
		}
	}
}
