package logger

import (
	"fmt"
	"os"
)

func Log(format string, a ...any) {
	fmt.Printf(format, a...)
	fmt.Println()
}

func ErrorExit(format string, a ...any) {
	fmt.Printf(format, a...)
	fmt.Println()
	os.Exit(1)
}
