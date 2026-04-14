package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func formatSize(size int64) string {
	const blockSize = 1024
	if size < blockSize {
		return fmt.Sprintf("%d B", size)
	}
	fSize := float64(size) / float64(blockSize)
	return fmt.Sprintf("%.1f Blks", fSize)
}

func getColor(info os.FileInfo) string {
	if info.IsDir() {
		return "1;36" // Cyan for directories
	}
	if info.Mode()&os.ModePerm&0111 != 0 {
		return "1;32" // Green for executables
	}
	ext := strings.ToLower(info.Name())
	if strings.HasSuffix(ext, ".go") || strings.HasSuffix(ext, ".py") || strings.HasSuffix(ext, ".sh") {
		return "1;33" // Yellow for source code
	}
	if strings.HasSuffix(ext, ".bpk") || strings.HasSuffix(ext, ".tar.gz") {
		return "1;35" // Magenta for packages/archives
	}
	return "0" // Default white
}

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	// Sort entries: directories first, then files
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})

	fmt.Printf("\r\n \033[1;36m[ BLOCK OS FILE EXPLORER ]\033[0m\r\n")
	fmt.Printf(" %-20s | %-12s | %-10s\r\n", "NAME", "SIZE", "TYPE")
	fmt.Println(" ---------------------------------------------------")

	for _, entry := range entries {
		info, _ := entry.Info()
		color := getColor(info)
		
		name := entry.Name()
		if entry.IsDir() {
			name = "[ " + name + " ]"
		}
		
		typeName := "FILE"
		if entry.IsDir() {
			typeName = "DIR"
		} else if info.Mode()&os.ModePerm&0111 != 0 {
			typeName = "EXE"
		}

		fmt.Printf(" \033[%sm%-20s\033[0m | %-12s | %s\r\n", 
			color, name, formatSize(info.Size()), typeName)
	}
	fmt.Println(" ---------------------------------------------------\r\n")
}
