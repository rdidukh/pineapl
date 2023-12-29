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

	result := ParseFile(tokens)

	log("")
	log("EXPRESSION PARSER RESULT")
	printExpressionResult(result)
	log("")

	if result.error != nil {
		errorExit("Expression parser error: %s", result.error)
	}
}

func printExpressionResult(result ExpressionParserResult, padding int) {
	paddingString := strings.Repeat(" ", 2 * padding)

	log("%s  %-2d %s size:%d ", paddingString, result.expressionType, result.size)
	for ()
}
