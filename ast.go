package main

import (
	"fmt"
	"maps"
)

type ExpressionType int

const (
	EXPRESSION_TYPE_UNKNOWN ExpressionType = iota
	EXPRESSION_TYPE_FUNCTION
)

func ParseExpressions(tokens []*Token) ([]ExpressionParserResult, error) {
	offset := 0
	results := []ExpressionParserResult{}

	parsers := oneOf(
		functionParser,
	)

	for offset < len(tokens) {
		log("   Iteration offset=%d len=%d", offset, len(tokens))
		result := parsers(tokens[offset:])

		results = append(results, result)

		if result.error != nil {
			log("     error result")
			return results, result.error
		}
		log("     ok result")
		offset += result.size
	}

	return results, nil
}

type ExpressionParserResult struct {
	offset         int
	size           int
	error          error
	expressionType ExpressionType
	tokenValues    map[string]*Token
}

type expressionParser func(tokens []*Token) ExpressionParserResult

type taggedExpressionParser struct {
	parser expressionParser
	tag    string
}

func requiredToken(tokenType TokenType) expressionParser {
	return requiredTokenWithValue(tokenType, "")
}

func requiredTokenWithValue(tokenType TokenType, key string) expressionParser {
	return func(tokens []*Token) ExpressionParserResult {
		if len(tokens) <= 0 {
			return ExpressionParserResult{
				error: fmt.Errorf("Expected %s, Found: EOF", tokenType),
			}
		}

		actualTokenType := tokens[0].tokenType

		if actualTokenType != tokenType {
			return ExpressionParserResult{
				error: fmt.Errorf("Expected %s, Found: %s", tokenType, actualTokenType),
			}
		}

		result := ExpressionParserResult{
			size: 1,
		}

		if len(key) > 0 {
			result.tokenValues = map[string]*Token{key: tokens[0]}
		}

		return result
	}
}

func optionalToken(expectedTokenType TokenType) expressionParser {
	return func(tokens []*Token) ExpressionParserResult {
		if len(tokens) <= 0 || tokens[0].tokenType != expectedTokenType {
			return ExpressionParserResult{}
		}

		return ExpressionParserResult{
			size: 1,
		}
	}
}

func allOf(parsers ...expressionParser) expressionParser {
	return func(tokens []*Token) ExpressionParserResult {
		result := ExpressionParserResult{
			tokenValues: map[string]*Token{},
		}
		for _, parser := range parsers {
			parserResult := parser(tokens[result.size:])
			maps.Copy(result.tokenValues, parserResult.tokenValues)
			result.size += parserResult.size

			if parserResult.error != nil {
				result.error = parserResult.error
				break
			}
		}

		return result
	}
}

func oneOf(parsers ...expressionParser) expressionParser {
	return func(tokens []*Token) ExpressionParserResult {
		bestResult := ExpressionParserResult{
			error: fmt.Errorf("Unexpected token: %s", tokens[0].tokenType),
		}

		for i, parser := range parsers {
			log("oneOf i = %d", i)
			result := parser(tokens)

			if result.error == nil {
				return result
			}

			if result.size > bestResult.size {
				bestResult = result
			}
		}

		return bestResult
	}
}

var functionParser expressionParser = func(tokens []*Token) ExpressionParserResult {
	const functionNameKey = "function.name"

	parser := allOf(
		requiredToken(TOKEN_TYPE_KEYWORD_FUNC),
		requiredToken(TOKEN_TYPE_WHITESPACE),
		requiredTokenWithValue(TOKEN_TYPE_IDENTIFIER, "function.name"),
		requiredToken(TOKEN_TYPE_ROUND_BRACKET_OPEN),
		optionalToken(TOKEN_TYPE_WHITESPACE),
		requiredToken(TOKEN_TYPE_ROUND_BRACKET_CLOSE),
		optionalToken(TOKEN_TYPE_WHITESPACE),
		requiredToken(TOKEN_TYPE_CURLY_BRACKET_OPEN),
		optionalToken(TOKEN_TYPE_WHITESPACE),
		requiredToken(TOKEN_TYPE_CURLY_BRACKET_CLOSE),
	)

	result := parser(tokens)

	result.expressionType = EXPRESSION_TYPE_FUNCTION

	return result
}

func (t ExpressionType) String() string {
	switch t {
	case EXPRESSION_TYPE_UNKNOWN:
		return "UNKNOWN"
	case EXPRESSION_TYPE_FUNCTION:
		return "FUNCTION"
	}

	panic(fmt.Sprintf("Unsupported expression type: %d", int(t)))
}
