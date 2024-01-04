package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

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

	output, err := compile(*inputFileFlag)

	if err != nil {
		log.Fatal(err)
	}

	logger.Log("")
	logger.Log("OUTPUT")
	logger.Log("%s", output)
	logger.Log("")
}

func compile(filename string) (string, error) {
	inputFileContents, err := os.ReadFile(filename)

	if err != nil {
		logger.ErrorExit("File read error: %s", err.Error())
	}

	file, err := ast.ParseString(string(inputFileContents))

	if err != nil {
		return "", fmt.Errorf("token error: %s", err)
	}

	printFile(file)

	if err != nil {
		return "", fmt.Errorf("expression parser error: %s", err)
	}

	output, err := file.Codegen()

	if err != nil {
		return "", fmt.Errorf("codegen error: %s", err)
	}

	return output, err
}

func printFile(file *ast.File) {
	jsonFile, err := json.MarshalIndent(file, "", " ")
	if err != nil {
		logger.ErrorExit("marshal error: %s", err)
	}

	logger.Log("")
	logger.Log("FILE")
	logger.Log(string(jsonFile))
	logger.Log("")
}
