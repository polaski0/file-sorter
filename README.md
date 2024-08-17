# File Sorter

A file sorter written in Go that makes use 
of goroutines to sort files faster which sort
files based on their file extensions.

## Usage

Create an executable file

```bash
make build
```

Run the executable file, supplying the `source` and `destination` directory

```bash
./fs target/directory output/directory
```

```bash
# Sort download directory and output the file to /tmp
./fs ~/Downloads /tmp
```
