package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
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

/*
varName123
(
)
_
123
123.45
"string"
=
+ - / * %
*/

type TokenType int

const (
	TOKEN_TYPE_UNKNOWN TokenType = iota
	TOKEN_TYPE_IDENTIFIER
	TOKEN_TYPE_WHITESPACE
	// TODO
	// TOKEN_TYPE_ROUND_BRACKET_OPEN
	// TOKEN_TYPE_ROUND_BRACKET_CLOSE
	// TOKEN_TYPE_EQUALS
	// TOKEN_TYPE_NUMBER
)

type Token struct {
	value     string
	tokenType TokenType
	start     int
	end       int
}

func (t TokenType) String() string {
	switch t {
	case TOKEN_TYPE_UNKNOWN:
		return "UNKNOWN"
	case TOKEN_TYPE_WHITESPACE:
		return "WHITESPACE"
	case TOKEN_TYPE_IDENTIFIER:
		return "IDENTIFIER"
	}

	panic("Unsupporte token type" + fmt.Sprint(int(t)))
}

func (t Token) String() string {
	return fmt.Sprintf("Token(type=%d, start=%d, end=%d, value=%q)", t.tokenType, t.start, t.end, t.value)
}

func getNextToken(sourceCode string, offset int) Token {
	log("getNextToken offset=%d, len=%d", offset, len(sourceCode))
	tokenGetters := []func(string, int) Token{
		getTokenWhitespace,
		getTokenIdentifier,
	}

	for _, tokenGetter := range tokenGetters {
		token := tokenGetter(sourceCode, offset)

		if token.tokenType != TOKEN_TYPE_UNKNOWN {
			return token
		}
	}

	return tokenTypeUnknown
}

func getTokenWhitespace(sourceCode string, offset int) Token {
	log("getTokenWhitespace %s %d\n", sourceCode, offset)

	pattern := regexp.MustCompile(`^\s+`)

	loc := pattern.FindStringIndex(sourceCode)

	return extractToken(sourceCode, offset, TOKEN_TYPE_WHITESPACE, loc)
}

func getTokenIdentifier(sourceCode string, offset int) Token {
	log("getTokenIdentifier %s %d", sourceCode, offset)

	pattern := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+`)

	loc := pattern.FindStringIndex(sourceCode)

	return extractToken(sourceCode, offset, TOKEN_TYPE_IDENTIFIER, loc)
}

func extractToken(sourceCode string, offset int, tokenType TokenType, loc []int) Token {
	if loc == nil {
		log("  loc == nil")
		return tokenTypeUnknown
	}

	start := loc[0]
	end := loc[1]
	token := sourceCode[start:end]

	if start != 0 {
		panic("start != 0")
	}

	log("  start=%d end=%d\n token=%s", start, end, token)
	log(string(sourceCode[loc[0]:loc[1]]))

	return Token{
		value:     token,
		tokenType: tokenType,
		start:     offset,
		end:       offset + end,
	}
}

func getTokens(code string) ([]Token, error) {
	log("getTokens")
	log(code)
	log("  len=%d\n", len(code))

	offset := 0
	tokens := []Token{}

	for offset < len(code) {
		token := getNextToken(code[offset:], offset)
		if token.tokenType == TOKEN_TYPE_UNKNOWN {
			return tokens, fmt.Errorf("Unknown token at offset %d: %q", offset, code[offset])
		}

		tokens = append(tokens, token)
		offset = token.end
	}

	return tokens, nil
}
