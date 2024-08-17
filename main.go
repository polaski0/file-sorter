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

	for _, f := range dir {
		if f.IsDir() {
			continue
		}

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

func ReadDir(path string) ([]fs.DirEntry, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	return dir, err
}

func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

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

func GetPathArgs() (target, output string) {
	args := os.Args
	target, output = "", "./tmp"

	if len(args) > 1 {
		if len(args) >= 2 {
			target = args[1]
		}

		if len(args) >= 3 {
			output = args[2]
		}
	} else {
		fmt.Println("Target directory is required")
		os.Exit(1)
	}

	return path.Clean(target), path.Clean(output)
}
