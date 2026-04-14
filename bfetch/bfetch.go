package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	Cyan   = "\033[1;36m"
	Blue   = "\033[1;34m"
	Yellow = "\033[1;33m"
	White  = "\033[1;37m"
	Reset  = "\033[0m"
	Bold   = "\033[1m"
)

func getUptime() string {
	data, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return "Unknown"
	}
	uptimeSeconds := strings.Fields(string(data))[0]
	d, _ := time.ParseDuration(uptimeSeconds + "s")
	return d.Round(time.Second).String()
}

func getMemory() string {
	data, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return "Unknown"
	}
	lines := strings.Split(string(data), "\n")
	var total, free uint64
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			fmt.Sscanf(line, "MemTotal: %d", &total)
		} else if strings.HasPrefix(line, "MemFree:") {
			fmt.Sscanf(line, "MemFree: %d", &free)
		}
	}
	used := (total - free) / 1024
	totalMb := total / 1024
	return fmt.Sprintf("%d MB / %d MB", used, totalMb)
}

func getKernel() string {
	data, err := ioutil.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(data))
}

func getPackageCount() int {
	files, err := ioutil.ReadDir("/repo/packages")
	if err != nil {
		return 0
	}
	return len(files)
}

func main() {
	logo := []string{
		Cyan + "    _______" + Reset,
		Cyan + "   /      /|" + Reset,
		Cyan + "  /______/ |" + Reset,
		Cyan + "  |      | |" + Reset,
		Cyan + "  |  " + White + Bold + "B" + Reset + Cyan + "   | |" + Reset,
		Cyan + "  |      | /" + Reset,
		Cyan + "  |______|/" + Reset,
	}

	host, _ := os.Hostname()
	user := os.Getenv("USER")
	if user == "" {
		user = "root"
	}

	info := []string{
		fmt.Sprintf("%s%s%s@%s%s%s", Blue, user, Reset, Blue, host, Reset),
		strings.Repeat("-", len(user)+len(host)+1),
		fmt.Sprintf("%sOS:%s      Block OS v2.0", Yellow, Reset),
		fmt.Sprintf("%sKernel:%s  %s", Yellow, Reset, getKernel()),
		fmt.Sprintf("%sUptime:%s  %s", Yellow, Reset, getUptime()),
		fmt.Sprintf("%sMemory:%s  %s", Yellow, Reset, getMemory()),
		fmt.Sprintf("%sPackages:%s %d (bpm)", Yellow, Reset, getPackageCount()),
		fmt.Sprintf("%sArch:%s     %s", Yellow, Reset, runtime.GOARCH),
	}

	fmt.Println()
	for i := 0; i < len(info) || i < len(logo); i++ {
		l := ""
		if i < len(logo) {
			l = logo[i]
		} else {
			l = "            "
		}

		r := ""
		if i < len(info) {
			r = info[i]
		}
		fmt.Printf("  %-25s %s\n", l, r)
	}
	fmt.Println()
}
