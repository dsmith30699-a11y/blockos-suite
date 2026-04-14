package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	Cyan   = "\033[1;36m"
	Yellow = "\033[1;33m"
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Green  = "\033[1;32m"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bcalc <expression> or <value> in <base>")
		fmt.Println("Example: bcalc \"1024 * 768\"")
		fmt.Println("Example: bcalc 255 in hex")
		return
	}

	input := strings.Join(os.Args[1:], " ")
	fmt.Printf("\n  %s%s[ BLOCK CALC ]%s - Dev Logic Engine\n", Cyan, Bold, Reset)
	fmt.Println("  --------------------------------------")

	// Check for "in hex/bin/dec" pattern
	if strings.Contains(input, " in ") {
		parts := strings.Split(input, " in ")
		valStr := strings.TrimSpace(parts[0])
		targetBase := strings.TrimSpace(parts[1])

		val, err := strconv.ParseInt(valStr, 0, 64)
		if err != nil {
			fmt.Printf("  Error: Could not parse value '%s'\n", valStr)
			return
		}

		fmt.Printf("  %sInput:%s    %d\n", Yellow, Reset, val)
		switch targetBase {
		case "hex":
			fmt.Printf("  %sHex:%s      %s0x%X%s\n", Yellow, Reset, Green, val, Reset)
		case "bin":
			fmt.Printf("  %sBinary:%s   %s0b%b%s\n", Yellow, Reset, Green, val, Reset)
		case "dec":
			fmt.Printf("  %sDecimal:%s  %s%d%s\n", Yellow, Reset, Green, val, Reset)
		default:
			fmt.Printf("  Error: Unknown base '%s'\n", targetBase)
		}
	} else {
		// Simple arithmetic placeholder (supports basic integers)
		// For a full evaluator, we'd need a parser, but let's handle the core bases.
		val, err := strconv.ParseInt(input, 0, 64)
		if err == nil {
			printAllBases(val)
		} else {
			fmt.Println("  Error: Simple arithmetic evaluator requires a full parser.")
			fmt.Println("  Try base conversion: bcalc 255 in hex")
		}
	}
	fmt.Println()
}

func printAllBases(val int64) {
	fmt.Printf("  %sDecimal:%s  %d\n", Yellow, Reset, val)
	fmt.Printf("  %sHex:%s      0x%X\n", Yellow, Reset, val)
	fmt.Printf("  %sBinary:%s   0b%b\n", Yellow, Reset, val)
	fmt.Printf("  %sOctal:%s    0%o\n", Yellow, Reset, val)
}
