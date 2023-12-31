package logger

import (
	"fmt"
	"os"
	"strings"
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

func LogPadded(padding int, format string, a ...any) {
	fmt.Print(strings.Repeat(" ", 2*padding))
	Log(format, a...)
}
