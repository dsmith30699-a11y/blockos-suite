package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

const (
	Cyan   = "\033[1;36m"
	Yellow = "\033[1;33m"
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Green  = "\033[1;32m"
)

type WriteCounter struct {
	Total      uint64
	ContentLen uint64
	StartTime  time.Time
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc *WriteCounter) PrintProgress() {
	width := 40
	percent := float64(wc.Total) / float64(wc.ContentLen) * 100
	completed := int(float64(width) * (float64(wc.Total) / float64(wc.ContentLen)))

	elapsed := time.Since(wc.StartTime).Seconds()
	speed := float64(wc.Total) / (1024 * 1024 * elapsed) // MB/s

	fmt.Printf("\r  %s[", Cyan)
	for i := 0; i < width; i++ {
		if i < completed {
			fmt.Print("█")
		} else {
			fmt.Print("░")
		}
	}
	fmt.Printf("]%s %s%5.1f%%%s | %s%.2f MB/s%s", Reset, Yellow, percent, Reset, Bold, speed, Reset)
}

func downloadFile(url string, filepath string) error {
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil // Always follow redirects
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	counter := &WriteCounter{
		ContentLen: uint64(resp.ContentLength),
		StartTime:  time.Now(),
	}

	fmt.Printf("\n  %s%s[ BLOCK GET ]%s - Secure Downloader\n", Cyan, Bold, Reset)
	fmt.Println("  --------------------------------------")
	fmt.Printf("  %sTarget:%s   %s\n", Yellow, Reset, filepath)
	fmt.Printf("  %sSize:%s     %.2f MB\n", Yellow, Reset, float64(resp.ContentLength)/(1024*1024))
	fmt.Println()

	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		return err
	}

	fmt.Print("\n\n")
	out.Close()
	err = os.Rename(filepath+".tmp", filepath)
	if err != nil {
		return err
	}

	fmt.Printf("  %s✔ Download Complete!%s\n\n", Green, Reset)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bget <url> [filename]")
		return
	}

	url := os.Args[1]
	var filename string

	if len(os.Args) > 2 {
		filename = os.Args[2]
	} else {
		filename = path.Base(url)
		if filename == "" || filename == "." || filename == "/" {
			filename = "downloaded_file"
		}
	}

	err := downloadFile(url, filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
