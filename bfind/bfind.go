package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	Cyan   = "\033[1;36m"
	Yellow = "\033[1;33m"
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Gray   = "\033[0;90m"
)

func search(root, pattern string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Skip inaccessible
		}
		if strings.Contains(strings.ToLower(d.Name()), strings.ToLower(pattern)) {
			results <- path
		}
		return nil
	})
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bfind <pattern> [directory]")
		return
	}

	pattern := os.Args[1]
	root := "/"
	if len(os.Args) > 2 {
		root = os.Args[2]
	}

	fmt.Printf("\n  %s%s[ BLOCK FIND ]%s - Parallel Search\n", Cyan, Bold, Reset)
	fmt.Println("  --------------------------------------")
	fmt.Printf("  %sSearching for:%s %s\n", Yellow, Reset, pattern)
	fmt.Printf("  %sIn Root:%s      %s\n", Yellow, Reset, root)
	fmt.Println()

	results := make(chan string, 100)
	var wg sync.WaitGroup
	var printerWg sync.WaitGroup

	printerWg.Add(1)
	foundCount := 0
	go func() {
		defer printerWg.Done()
		for res := range results {
			foundCount++
			// Highlight the pattern in the path
			displayPath := strings.ReplaceAll(res, pattern, Cyan+pattern+Reset)
			fmt.Printf("  %s-> %s%s\n", Gray, Reset, displayPath)
		}
	}()

	// Get top-level directories to parallelize
	entries, _ := os.ReadDir(root)
	for _, entry := range entries {
		wg.Add(1)
		fullPath := filepath.Join(root, entry.Name())
		go search(fullPath, pattern, &wg, results)
	}

	wg.Wait()
	close(results)
	printerWg.Wait()

	fmt.Printf("\n  %s✔ Search Complete. Found %d matches.%s\n\n", Cyan, foundCount, Reset)
}
