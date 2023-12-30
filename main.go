package main

import (
	"flag"
	"os"
	"strings"

	"github.com/rdidukh/pineapl/ast"
	"github.com/rdidukh/pineapl/logger"
)

var inputFileFlag = flag.String("i", "", "TODO: usage")

func main() {
	logger.Log("main [start]")

	flag.Parse()

	if *inputFileFlag == "" {
		logger.ErrorExit("Missing -i")
	}

	inputFileContents, err := os.ReadFile(*inputFileFlag)

	if err != nil {
		logger.ErrorExit("File read error: %s", err.Error())
	}

	tokens, file, err := ast.ParseString(string(inputFileContents))

	logger.Log("")
	logger.Log("TOKENS")
	for i, t := range tokens {
		logger.Log("  %-3d %-20s %3d..%2d   %-10q", i, t.Type, t.Start, t.End, t.Value)
	}
	logger.Log("")

	if err != nil {
		logger.ErrorExit("Token error: %s", err)
	}

	printFile(file)

	if err != nil {
		logger.ErrorExit("Expression parser error: %s", err)
	}
}

func printFile(file *ast.File) {
	logger.Log("")
	logger.Log("FILE ")
	for _, function := range file.Functions {
		printFunction(function, 1)
	}
	logger.Log("")
}

func printFunction(function *ast.Function, padding int) {
	paddingString := strings.Repeat(" ", 2*padding)
	logger.Log("%sFUNCTION %s", paddingString, function.Name)
	for _, parameter := range function.Parameters {
		printParameter(parameter, padding+1)
	}
}

func printParameter(parameter *ast.Parameter, padding int) {
	paddingString := strings.Repeat(" ", 2*padding)
	logger.Log("%sPARAMETER %s %s", paddingString, parameter.Name, parameter.Type)
}
