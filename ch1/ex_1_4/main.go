package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	filenames := make(map[string]string)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, "-", counts, filenames)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, arg, counts, filenames)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%s%d\t%s\n\n", filenames[line], n, line)
		}
	}
}

func countLines(f *os.File, filename string, counts map[string]int, filenames map[string]string) {
	localCounts := make(map[string]int)
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := input.Text()
		counts[line]++
		localCounts[line]++
		if localCounts[line] == 1 {
			filenames[line] += filename + "\n"
		}
	}
	// NOTE: ignoring potential errors from input.Err()
}
