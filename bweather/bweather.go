package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	Cyan   = "\033[1;36m"
	Yellow = "\033[1;33m"
	Reset  = "\033[0m"
	Bold   = "\033[1m"
)

func main() {
	location := ""
	if len(os.Args) > 1 {
		location = strings.Join(os.Args[1:], "+")
	}

	fmt.Printf("\n  %s%s[ BLOCK WEATHER v2.0 ]%s\n", Cyan, Bold, Reset)
	fmt.Println("  -----------------------------------------")
	if location == "" {
		fmt.Println("  Auto-detecting local weather...")
	} else {
		fmt.Printf("  Fetching weather for: %s%s%s\n", Yellow, strings.ReplaceAll(location, "+", " "), Reset)
	}
	fmt.Println()

	// wttr.in format: ?0=current, ?n=narrow, ?F=no follow (use curl agent for ASCII)
	url := "https://wttr.in/" + location + "?0&n"
	
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	
	// CRITICAL: wttr.in uses User-Agent to decide between HTML and ASCII
	req.Header.Set("User-Agent", "curl/7.79.1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("  Error: Could not reach weather service: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("  Error: Could not read response: %v\n", err)
		os.Exit(1)
	}

	// Print the result with a small indent
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			fmt.Printf("  %s\n", line)
		}
	}

	fmt.Println("\n  [ DATA PROVIDED BY WTTR.IN ]")
}
