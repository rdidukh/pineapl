package main

import (
	"flag"
	"fmt"
	"os"
)

var inputFileFlag = flag.String("i", "", "TODO: usage")

var tokenTypeUnknown = Token{
	tokenType: TOKEN_TYPE_UNKNOWN,
}

func log(format string, a ...any) {
	fmt.Printf(format, a...)
	fmt.Println()
}

func errorExit(format string, a ...any) {
	fmt.Printf(format, a...)
	fmt.Println()
	os.Exit(1)
}

func main() {
	log("main [start]")

	flag.Parse()

	log("  after flag.Parse()")

	if *inputFileFlag == "" {
		errorExit("Missing -i")
	}

	inputFileContents, err := os.ReadFile(*inputFileFlag)

	log("  after ReadFile")

	if err != nil {
		errorExit("File read error: %s", err.Error())
	}

	codeToCompile := string(inputFileContents)

	log("  codeToCompile len=%d", len(codeToCompile))

	tokens, err := getTokens(codeToCompile)

	log("After getTokens, len=%d", len(tokens))
	for _, token := range tokens {
		log("  %s", token)
	}

	if err != nil {
		errorExit("Token error: %s", err)
	}
}
