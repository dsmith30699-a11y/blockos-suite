package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
)

const (
	Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	Cyan    = "\033[1;36m"
	Yellow  = "\033[1;33m"
	Red     = "\033[1;31m"
	Reset   = "\033[0m"
	Bold    = "\033[1m"
)

func generatePassword(length int) (string, error) {
	password := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(Charset))))
		if err != nil {
			return "", err
		}
		password[i] = Charset[num.Int64()]
	}
	return string(password), nil
}

func main() {
	length := 16
	if len(os.Args) > 1 {
		l, err := strconv.Atoi(os.Args[1])
		if err == nil && l > 0 {
			length = l
		}
	}

	password, err := generatePassword(length)
	if err != nil {
		fmt.Printf("%sError generating password: %v%s\n", Red, err, Reset)
		return
	}

	fmt.Println()
	fmt.Printf("  %s%s[ BLOCK PASS ]%s - Secure Generator\n", Cyan, Bold, Reset)
	fmt.Println("  --------------------------------------")
	fmt.Printf("  %sLength:%s   %d characters\n", Yellow, Reset, length)
	fmt.Printf("  %sEntropy:%s  High (Crypto/Rand)\n", Yellow, Reset)
	fmt.Println()
	fmt.Printf("  %sGenerated Password:%s\n", Bold, Reset)
	fmt.Printf("  %s%s%s\n", Cyan, password, Reset)
	fmt.Println()
	fmt.Printf("  %s[ KEEP THIS SECURE ]%s\n", Red, Reset)
	fmt.Println()
}
