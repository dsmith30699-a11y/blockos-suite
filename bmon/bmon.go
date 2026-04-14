package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Process struct {
	PID     int
	Command string
	CPU     float64
	Mem     float64
}

const (
	Cyan   = "\033[1;36m"
	Yellow = "\033[1;33m"
	Green  = "\033[1;32m"
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Gray   = "\033[0;90m"
)

func getSystemStats() (float64, float64) {
	// CPU
	data, _ := ioutil.ReadFile("/proc/stat")
	fields := strings.Fields(strings.Split(string(data), "\n")[0])[1:]
	var total uint64
	for _, f := range fields {
		v, _ := strconv.ParseUint(f, 10, 64)
		total += v
	}
	idle, _ := strconv.ParseUint(fields[3], 10, 64)
	cpu := 100.0 * (1.0 - float64(idle)/float64(total))

	// Mem
	data, _ = ioutil.ReadFile("/proc/meminfo")
	var t, a uint64
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			fmt.Sscanf(line, "MemTotal: %d", &t)
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fmt.Sscanf(line, "MemAvailable: %d", &a)
		}
	}
	mem := 100.0 * (1.0 - float64(a)/float64(t))
	return cpu, mem
}

func getTopProcesses() []Process {
	dirs, _ := ioutil.ReadDir("/proc")
	var procs []Process
	for _, d := range dirs {
		pid, err := strconv.Atoi(d.Name())
		if err != nil { continue }

		cmd, _ := ioutil.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
		stat, _ := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
		fields := strings.Fields(string(stat))
		if len(fields) < 15 { continue }
		
		utime, _ := strconv.ParseFloat(fields[13], 64)
		stime, _ := strconv.ParseFloat(fields[14], 64)
		
		procs = append(procs, Process{
			PID:     pid,
			Command: strings.TrimSpace(string(cmd)),
			CPU:     (utime + stime) / 100.0, // Simplified CPU metric
			Mem:     0.1, // Placeholder for minimal env
		})
	}
	sort.Slice(procs, func(i, j int) bool { return procs[i].CPU > procs[j].CPU })
	if len(procs) > 5 { return procs[:5] }
	return procs
}

func drawBar(label string, percent float64, color string) {
	width := 25
	filled := int(percent / 100.0 * float64(width))
	if filled > width { filled = width }
	fmt.Printf("  %-5s %s[", label, Bold)
	for i := 0; i < width; i++ {
		if i < filled { fmt.Print(color + "█" + Reset) } else { fmt.Print(Gray + "░" + Reset) }
	}
	fmt.Printf("%s] %5.1f%%\n", Bold, percent)
}

func main() {
	fmt.Print("\033[?1049h\033[H\033[2J") // Alternate buffer
	defer fmt.Print("\033[?1049l")

	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		cpu, mem := getSystemStats()
		procs := getTopProcesses()

		fmt.Print("\033[H")
		fmt.Printf("\n  %s%s[ BLOCK MONITOR v1.0 ]%s\n", Cyan, Bold, Reset)
		fmt.Println("  --------------------------------------")
		drawBar("CPU", cpu, Green)
		drawBar("RAM", mem, Yellow)
		fmt.Println("  --------------------------------------")
		fmt.Printf("  %s%-6s %-12s %-6s%s\n", Bold, "PID", "COMMAND", "USAGE", Reset)
		
		for _, p := range procs {
			cmd := p.Command
			if len(cmd) > 12 { cmd = cmd[:11] + "+" }
			fmt.Printf("  %-6d %-12s %-6.1f\n", p.PID, cmd, p.CPU)
		}
		fmt.Printf("\n  %s[ Ctrl+C to Exit ]%s", Gray, Reset)
	}
}
