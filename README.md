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
./fs [src/directory]... dest/directory
```

```bash
# Sorts "Downloads" directory to "/tmp" directory
./fs ~/Downloads /tmp

# Sorts both "Downloads" and "Documents" directory to "/tmp" directory
./fs ~/Downloads ~/Documents /tmp
```
