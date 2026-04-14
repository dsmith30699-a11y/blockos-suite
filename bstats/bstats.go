package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func getCPUUsage() float64 {
	file, err := os.Open("/proc/stat")
	if err != nil { return 0.0 }
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	fields := strings.Fields(scanner.Text())[1:]
	var total uint64
	for _, f := range fields {
		v, _ := strconv.ParseUint(f, 10, 64)
		total += v
	}
	idle, _ := strconv.ParseUint(fields[3], 10, 64)
	return 100.0 * (1.0 - float64(idle)/float64(total))
}

func getMemUsage() float64 {
	file, err := os.Open("/proc/meminfo")
	if err != nil { return 0.0 }
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var total, available uint64
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fmt.Sscanf(line, "MemTotal: %d", &total)
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fmt.Sscanf(line, "MemAvailable: %d", &available)
		}
	}
	if total == 0 { return 0.0 }
	return 100.0 * (1.0 - float64(available)/float64(total))
}

func drawBar(label string, percent float64, color string) {
	width := 30
	filled := int(percent / 100.0 * float64(width))
	fmt.Printf(" \033[0m%-5s [\033[%sm", label, color)
	for i := 0; i < width; i++ {
		if i < filled {
			fmt.Print("█")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Printf("\033[0m] %5.1f%%\n", percent)
}

func main() {
	fmt.Print("\033[?1049h\033[H\033[2J") // Alternate buffer & Clear
	defer fmt.Print("\033[?1049l")

	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		fmt.Print("\033[H")
		fmt.Println("\r\n \033[1;36m[ BLOCK OS REAL-TIME MONITOR ]\033[0m")
		fmt.Println(" ----------------------------------------")
		drawBar("CPU", getCPUUsage(), "1;32") // Green
		drawBar("RAM", getMemUsage(), "1;33") // Yellow
		fmt.Println(" ----------------------------------------")
		fmt.Println("\r\n Press Ctrl+C to Exit")
	}
}
