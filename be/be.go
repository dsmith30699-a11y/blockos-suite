package main

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

var (
	rows     []string
	cursorX  = 0
	cursorY  = 0
	filename string
	screenH  = 24
	screenW  = 80
	status   string
	statusT  time.Time
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: be <filename>")
		return
	}
	filename = os.Args[1]
	loadFile(filename)

	oldState := setRawMode()
	defer restoreMode(oldState)

	fmt.Print("\033[?1049h") // Alternate buffer
	defer fmt.Print("\033[?1049l")

	for {
		getTermSize()
		draw()
		buf := make([]byte, 8)
		n, _ := os.Stdin.Read(buf)
		if n == 0 { continue }

		if buf[0] == 17 { break } // Ctrl+Q
		if buf[0] == 19 { 
			if err := saveFile(); err == nil {
				setStatus("[ File Saved! ]")
			} else {
				setStatus("[ Error Saving File! ]")
			}
			continue 
		}
		handleInput(buf[:n])
	}
}

func setStatus(msg string) {
	status = msg
	statusT = time.Now()
}

func getTermSize() {
	var ws struct { Row, Col, X, Y uint16 }
	syscall.Syscall(syscall.SYS_IOCTL, os.Stdout.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&ws)))
	screenH = int(ws.Row)
	screenW = int(ws.Col)
}

func draw() {
	fmt.Print("\033[H")
	fmt.Print("\033[1;36m[ BLOCK OS EDITOR - BE v0.4 ]\033[0m\r\n")
	
	for i := 0; i < screenH-3; i++ {
		if i < len(rows) {
			fmt.Printf("\033[K%s\r\n", rows[i])
		} else {
			fmt.Print("\033[K~\r\n")
		}
	}
	
	// Status Message (clears after 3 seconds)
	currentStatus := ""
	if time.Since(statusT) < 3*time.Second {
		currentStatus = status
	}

	// Status Bar
	fmt.Printf("\033[%d;1H\033[K\033[1;33m %s\033[0m", screenH-1, currentStatus)
	fmt.Printf("\033[%d;1H\033[7m FILE: %s | CTRL+S: Save | CTRL+Q: Exit \033[0m", screenH, filename)
	fmt.Printf("\033[%d;%dH", cursorY+2, cursorX+1)
}

func handleInput(b []byte) {
	if b[0] == 27 && len(b) >= 3 {
		switch b[2] {
		case 65: if cursorY > 0 { cursorY-- }
		case 66: if cursorY < len(rows)-1 { cursorY++ }
		case 67: if cursorX < len(rows[cursorY]) { cursorX++ }
		case 68: if cursorX > 0 { cursorX-- }
		}
		return
	}
	if b[0] == 127 || b[0] == 8 {
		if cursorX > 0 {
			rows[cursorY] = rows[cursorY][:cursorX-1] + rows[cursorY][cursorX:]
			cursorX--
		}
		return
	}
	if b[0] == 13 {
		rows = append(rows[:cursorY+1], append([]string{""}, rows[cursorY+1:]...)...)
		cursorY++; cursorX = 0
		return
	}
	if b[0] >= 32 && b[0] <= 126 {
		rows[cursorY] = rows[cursorY][:cursorX] + string(b[0]) + rows[cursorY][cursorX:]
		cursorX++
	}
}

func loadFile(fn string) {
	file, err := os.Open(fn)
	if err != nil { rows = []string{""}; return }
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { rows = append(rows, scanner.Text()) }
	if len(rows) == 0 { rows = []string{""} }
}

func saveFile() error {
	file, err := os.Create(filename)
	if err != nil { return err }
	defer file.Close()
	for _, line := range rows { fmt.Fprintln(file, line) }
	return nil
}

func setRawMode() *syscall.Termios {
	var old syscall.Termios
	syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), syscall.TCGETS, uintptr(unsafe.Pointer(&old)))
	new := old
	new.Lflag &^= syscall.ICANON | syscall.ECHO
	new.Iflag &^= syscall.IXON 
	syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), syscall.TCSETS, uintptr(unsafe.Pointer(&new)))
	return &old
}

func restoreMode(old *syscall.Termios) {
	if old != nil { syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), syscall.TCSETS, uintptr(unsafe.Pointer(old))) }
}
