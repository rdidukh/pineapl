package main

import (
	"fmt"
)

type ExpressionType int

const (
	EXPRESSION_TYPE_UNKNOWN ExpressionType = iota
	EXPRESSION_TYPE_GROUP
	EXPRESSION_TYPE_FILE
	EXPRESSION_TYPE_FUNCTION
)

type expressionParserRequest struct {
	tokens []*Token
}

type ExpressionParserResult struct {
	size               int
	error              error
	expressionType     ExpressionType
	tokenValue         *Token
	intValues          []int
	childResults       []ExpressionParserResult
	taggedChildResults map[string]ExpressionParserResult
}

func (r ExpressionParserResult) addIntValues(values ...int) {
	r.intValues = append(r.intValues, values...)
}

func (r ExpressionParserResult) addChildResult(result ...ExpressionParserResult) {
	r.childResults = append(r.childResults, result...)
}

func (r ExpressionParserResult) addTaggedChildResult(result ExpressionParserResult, tag string) {
	if tag == "" {
		return
	}
	if r.taggedChildResults == nil {
		r.taggedChildResults = map[string]ExpressionParserResult{}
	}
	r.taggedChildResults[tag] = result
}

func (r ExpressionParserResult) addTaggedChildResults(results map[string]ExpressionParserResult) {
	if r.taggedChildResults == nil {
		r.taggedChildResults = map[string]ExpressionParserResult{}
	}
	for tag, result := range results {
		r.taggedChildResults[tag] = result
	}
}

type expressionParser func(expressionParserRequest) ExpressionParserResult
type expressionParserCallback func(result ExpressionParserResult)

type expressionParserConfig struct {
	parser expressionParser
	tag    string
}

func ParseFile(tokens []*Token) ExpressionParserResult {
	fileParser := oneOfRepeated(expressionParserConfig{
		parser: functionParser,
		tag:    "function",
	})

	result := fileParser(expressionParserRequest{
		tokens: tokens,
	})

	result.expressionType = EXPRESSION_TYPE_FILE

	return result
}

func requiredToken(tokenType TokenType) expressionParserConfig {
	return requiredTokenWithTag(tokenType, "")
}

func requiredTokenWithTag(tokenType TokenType, tag string) expressionParserConfig {
	return expressionParserConfig{
		parser: requiredTokenParser(tokenType),
		tag:    tag,
	}
}

func requiredTokenParser(tokenType TokenType) expressionParser {
	return func(request expressionParserRequest) ExpressionParserResult {
		tokens := request.tokens
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
			tokenValue: tokens[0],
			size:       1,
		}
	}
}

func optionalToken(tokenType TokenType) expressionParserConfig {
	return expressionParserConfig{
		parser: optionalTokenParser(tokenType),
	}
}

func optionalTokenParser(tokenType TokenType) expressionParser {
	return func(request expressionParserRequest) ExpressionParserResult {
		tokens := request.tokens
		if len(tokens) <= 0 || tokens[0].tokenType != tokenType {
			return ExpressionParserResult{}
		}

		return ExpressionParserResult{
			size:       1,
			tokenValue: tokens[0],
		}
	}
}

func oneOfRepeated(configs ...expressionParserConfig) expressionParser {
	return func(request expressionParserRequest) ExpressionParserResult {
		tokens := request.tokens
		result := ExpressionParserResult{
			expressionType: EXPRESSION_TYPE_GROUP,
		}

		parsers := []expressionParser{}

		for _, config := range configs {
			parsers = append(parsers, config.parser)
		}

		parser := oneOf(configs...)

		for result.size < len(tokens) {
			log("   Iteration offset=%d len=%d", result.size, len(tokens))

			childResult := parser(expressionParserRequest{
				tokens: tokens[result.size:],
			})

			result.addChildResult(childResult)
			result.addTaggedChildResults(childResult.taggedChildResults)
			result.error = childResult.error
			result.size += childResult.size
			result.addIntValues(childResult.intValues...)

			if childResult.error != nil {
				log("     error result")
				return result
			}
			log("     ok result")
		}

		return result
	}
}

func allOfOrdered(configs ...expressionParserConfig) expressionParser {
	return func(request expressionParserRequest) ExpressionParserResult {
		result := ExpressionParserResult{
			taggedChildResults: map[string]ExpressionParserResult{},
			expressionType:     EXPRESSION_TYPE_GROUP,
		}
		for _, config := range configs {
			parserRequest := request
			parserRequest.tokens = parserRequest.tokens[result.size:]
			childResult := config.parser(parserRequest)

			result.addChildResult(childResult)
			result.addTaggedChildResult(childResult, config.tag)
			result.size += childResult.size

			if childResult.error != nil {
				result.error = childResult.error
				break
			}
		}

		return result
	}
}

func oneOf(configs ...expressionParserConfig) expressionParser {
	return func(request expressionParserRequest) ExpressionParserResult {
		tokens := request.tokens
		bestResult := ExpressionParserResult{
			error: fmt.Errorf("Unexpected token: %s", tokens[0].tokenType),
		}

		for i, config := range configs {
			log("oneOf i = %d", i)
			childResult := config.parser(request)

			result := ExpressionParserResult{}
			result.addChildResult(childResult)
			result.addTaggedChildResult(childResult, config.tag)
			result.error = childResult.error
			result.addIntValues(i)
			result.expressionType = EXPRESSION_TYPE_GROUP

			if result.error == nil {
				return result
			}

			if result.size >= bestResult.size {
				bestResult = result
			}
		}

		return bestResult
	}
}

func functionParser(request expressionParserRequest) ExpressionParserResult {
	const functionNameTag = "function.name"

	parser := allOfOrdered(
		requiredToken(TOKEN_TYPE_KEYWORD_FUNC),
		requiredToken(TOKEN_TYPE_WHITESPACE),
		requiredTokenWithTag(TOKEN_TYPE_IDENTIFIER, functionNameTag),
		requiredToken(TOKEN_TYPE_ROUND_BRACKET_OPEN),
		optionalToken(TOKEN_TYPE_WHITESPACE),
		requiredToken(TOKEN_TYPE_ROUND_BRACKET_CLOSE),
		optionalToken(TOKEN_TYPE_WHITESPACE),
		requiredToken(TOKEN_TYPE_CURLY_BRACKET_OPEN),
		optionalToken(TOKEN_TYPE_WHITESPACE),
		requiredToken(TOKEN_TYPE_CURLY_BRACKET_CLOSE),
	)

	result := parser(request)

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
