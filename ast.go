package main

import "fmt"

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
	offset int
	size   int
	error  error
}

type expressionParser func(tokens []*Token) ExpressionParserResult

func requiredToken(tokenType TokenType) expressionParser {
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

		return ExpressionParserResult{
			size: 1,
		}
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
		offset := 0
		for _, parser := range parsers {
			result := parser(tokens[offset:])
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

var functionParser = allOf(
	requiredToken(TOKEN_TYPE_KEYWORD_FUNC),
	requiredToken(TOKEN_TYPE_WHITESPACE),
	requiredToken(TOKEN_TYPE_IDENTIFIER),
	requiredToken(TOKEN_TYPE_ROUND_BRACKET_OPEN),
	optionalToken(TOKEN_TYPE_WHITESPACE),
	requiredToken(TOKEN_TYPE_ROUND_BRACKET_CLOSE),
	optionalToken(TOKEN_TYPE_WHITESPACE),
	requiredToken(TOKEN_TYPE_CURLY_BRACKET_OPEN),
	optionalToken(TOKEN_TYPE_WHITESPACE),
	requiredToken(TOKEN_TYPE_CURLY_BRACKET_CLOSE),
)
