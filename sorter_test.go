package main

import (
	"fmt"
	"os"
	"path"
	"slices"
	"sync"
	"testing"
)

const dir = "./tmp"

// Create test files on the current directory
func setup() {
	var wg sync.WaitGroup
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
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := os.WriteFile(path.Join(dest, f.filename), []byte(f.value), 0644)
				if err != nil {
					fmt.Printf("Error writing file %v, %v", f.filename, err)
					return
				}
			}()
		}
	}

	wg.Wait()
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
	var wg sync.WaitGroup
	setup()
	defer teardown()

	dest := path.Join(dir, "out")
	src := []string{
		path.Join(dir, "dir1"),
		path.Join(dir, "dir2"),
	}

	testCases := map[string][]string{
		"md": []string{
			"README.md",
			"Fizz.md",
			"README (1).md",
		},
		"txt": []string{
			"test.txt",
			"test (1).txt",
			"test (1) (1).txt",
			"test (2).txt",
			"test (3).txt",
		},
	}

	sr := NewSorter(src, dest)
	err := sr.Start()
	if err != nil {
		t.Errorf("Error %v\n", err)
	}

	for k, v := range testCases {
		wg.Add(1)

		go func() {
			defer wg.Done()
			fsDir, err := os.ReadDir(path.Join(dest, k))
			if err != nil {
				t.Errorf("Error %v", err)
			}
			files := []string{}

			for _, file := range fsDir {
				files = append(files, file.Name())
			}

			for _, file := range v {
				ok := slices.Contains(files, file)
				if !ok {
					t.Errorf("File %v not found", file)
				}
			}
		}()
	}

	wg.Wait()
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

	for _, val := range testCases {
		exists, err := isFileExists(val.file, val.dir)
		if err != nil && !os.IsNotExist(err) {
			t.Errorf("Error %v", err)
		}

		if exists != val.expected {
			t.Errorf("File: %v; Found %v, expected %v\n",
				val.file,
				exists,
				val.expected)
		}
	}
}
