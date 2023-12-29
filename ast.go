package main

import "fmt"

func ParseExpressions(tokens []*Token) ([]ExpressionParserResult, error) {
	offset := 0
	results := []ExpressionParserResult{}

	parsers := oneOf([]expressionParser{
		functionParser,
	})

	for offset < len(tokens) {
		log("   Iteration offset=%d len=%d", offset, len(tokens))
		result := parsers(tokens[offset:], offset)

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

type tokenCount int

const (
	TOKEN_COUNT_UNKNOWN tokenCount = iota
	TOKEN_COUNT_ZERO_OR_ONE
	TOKEN_COUNT_ONE
)

type ExpressionParserResult struct {
	offset int
	size   int
	error  error
}

type expressionParser func(tokens []*Token, offset int) ExpressionParserResult

func requiredToken(tokenType TokenType) expressionParser {
	return func(tokens []*Token, offset int) ExpressionParserResult {
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

		return ExpressionParserResult{
			offset: offset,
			size:   1,
		}
	}
}

func optionalToken(expectedTokenType TokenType) expressionParser {
	return func(tokens []*Token, offset int) ExpressionParserResult {
		if len(tokens) <= 0 {
			return ExpressionParserResult{
				// TODO: offset: offset,
				size: 0,
			}
		}

		if tokens[0].tokenType != expectedTokenType {
			return ExpressionParserResult{
				// TODO: offset: offset,
				size: 0,
			}
		}

		return ExpressionParserResult{
			// TODO: offset: offset,
			size: 1,
		}
	}
}

func group(parsers []expressionParser) expressionParser {
	return func(tokens []*Token, offset int) ExpressionParserResult {
		offset = 0
		for _, parser := range parsers {
			// func(tokens []*Token, offset int) parseResult
			result := parser(tokens[offset:], offset)
			offset += result.size

			if result.error != nil {
				return ExpressionParserResult{
					size:  offset,
					error: result.error,
				}
			}
		}

		return ExpressionParserResult{
			size: offset,
		}
	}
}

func oneOf(parsers []expressionParser) expressionParser {
	return func(tokens []*Token, offset int) ExpressionParserResult {
		bestResult := ExpressionParserResult{
			error: fmt.Errorf("Unexpected token: %s", tokens[0].tokenType),
		}

		for i, parser := range parsers {
			log("oneOf i = %d", i)
			result := parser(tokens, offset)

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

var functionParser = group([]expressionParser{
	requiredToken(TOKEN_TYPE_KEYWORD_FUNC),
	requiredToken(TOKEN_TYPE_WHITESPACE),
	requiredToken(TOKEN_TYPE_IDENTIFIER),
	requiredToken(TOKEN_TYPE_ROUND_BRACKET_OPEN),
	optionalToken(TOKEN_TYPE_WHITESPACE),
	requiredToken(TOKEN_TYPE_ROUND_BRACKET_CLOSE),
	optionalToken(TOKEN_TYPE_WHITESPACE),
	requiredToken(TOKEN_TYPE_CURLY_BRACKET_OPEN),
	optionalToken(TOKEN_TYPE_WHITESPACE),
	requiredToken(TOKEN_TYPE_CURLY_BRACKET_CLOSE)},
)
