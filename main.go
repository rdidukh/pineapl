package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var inputFileFlag = flag.String("i", "", "TODO: usage")

var tokenTypeUnknown = &Token{
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

	tokens, err := GetTokens(codeToCompile)

	log("")
	log("TOKENS")
	for i, token := range tokens {
		log("  %-3d %-20s %3d..%2d   %-10q", i, token.tokenType, token.start, token.end, token.value)
	}
	log("")

	if err != nil {
		errorExit("Token error: %s", err)
	}

	file, err := ParseFile(tokens)

	printFile(file)

	if err != nil {
		errorExit("Expression parser error: %s", err)
	}
}

func printFile(file *File) {
	log("")
	log("FILE ")
	for _, function := range file.functions {
		printFunction(function, 1)
	}
	log("")
}

func printFunction(function *Function, padding int) {
	paddingString := strings.Repeat(" ", 2*padding)
	log("%sFUNCTION %s", paddingString, function.name)
	for _, parameter := range function.parameters {
		printParameter(parameter, padding+1)
	}
}

func printParameter(parameter *Parameter, padding int) {
	paddingString := strings.Repeat(" ", 2*padding)
	log("%sPARAMETER %s %s", paddingString, parameter.name, parameter.paramType)
}
