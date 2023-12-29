package main

import (
	"fmt"
)

type Expression struct {
	token     *Token
	function  *Function
	file      *File
	parameter *Parameter
}

type Function struct {
	name       string
	parameters []*Parameter
}

type Parameter struct {
	name      string
	paramType string
}

type File struct {
	functions []*Function
}

type ExpressionType int

const (
	EXPRESSION_TYPE_UNKNOWN ExpressionType = iota
	EXPRESSION_TYPE_GROUP
	EXPRESSION_TYPE_FILE
	EXPRESSION_TYPE_FUNCTION
)

type expressionParserRequest struct {
	tokens   []*Token
	callback expressionParserCallback
}

type ExpressionParserResult struct {
	size       int
	error      error
	expression *Expression
}

type expressionParser func(expressionParserRequest) ExpressionParserResult
type expressionParserCallback func(result ExpressionParserResult)

type expressionParserConfig struct {
	parser   expressionParser
	callback expressionParserCallback
}

func (c expressionParserConfig) onSuccess(result ExpressionParserResult) {
	if c.callback != nil {
		c.callback(result)
	}
}

func ParseFile(tokens []*Token) (*File, error) {
	file := &File{}

	_, err := parseOneOfRepeated(expressionParserRequest{tokens: tokens},
		expressionParserConfig{
			parser: functionParser,
			callback: func(result ExpressionParserResult) {
				file.functions = append(file.functions, result.expression.function)
			},
		})

	return file, err
}

func requiredToken(tokenType TokenType) expressionParserConfig {
	return requiredTokenWithCallback(tokenType, func(result ExpressionParserResult) {})
}

func requiredTokenWithCallback(tokenType TokenType, callback expressionParserCallback) expressionParserConfig {
	return expressionParserConfig{
		parser:   requiredTokenParser(tokenType),
		callback: callback,
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
			expression: &Expression{token: tokens[0]},
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
			expression: &Expression{token: tokens[0]},
			size:       1,
		}
	}
}

func oneOfRepeatedParser(configs ...expressionParserConfig) expressionParser {
	return func(request expressionParserRequest) ExpressionParserResult {
		size, err := parseOneOfRepeated(request, configs...)
		return ExpressionParserResult{
			size:  size,
			error: err,
		}
	}
}

func parseOneOfRepeated(request expressionParserRequest, configs ...expressionParserConfig) (int, error) {
	return parseOneOfRepeatedUntil(request, TOKEN_TYPE_UNKNOWN, configs...)
}

func oneOfRepeatedUntilParser(terminator TokenType, configs ...expressionParserConfig) expressionParser {
	return func(request expressionParserRequest) ExpressionParserResult {
		size, err := parseOneOfRepeatedUntil(request, terminator, configs...)
		return ExpressionParserResult{
			size:  size,
			error: err,
		}
	}
}

func parseOneOfRepeatedUntil(request expressionParserRequest, terminator TokenType, configs ...expressionParserConfig) (int, error) {
	offset := 0
	tokens := request.tokens

	for offset < len(tokens) && tokens[offset].tokenType != terminator {
		size, err := parseOneOf(expressionParserRequest{
			tokens: request.tokens[offset:],
		}, configs...)

		if err != nil {
			return size, err
		}

		offset += size
	}

	return offset, nil
}

func allOrderedParser(configs ...expressionParserConfig) expressionParser {
	return func(request expressionParserRequest) ExpressionParserResult {
		size, err := parseAllOrdered(request, configs...)
		return ExpressionParserResult{
			size:  size,
			error: err,
		}
	}
}

func parseAllOrdered(request expressionParserRequest, configs ...expressionParserConfig) (int, error) {
	offset := 0
	for _, config := range configs {
		result := config.parser(expressionParserRequest{
			tokens: request.tokens[offset:],
		})

		offset += result.size

		if result.error != nil {
			return offset, result.error
		}

		config.onSuccess(result)
	}

	return offset, nil
}

func oneOfParser(configs ...expressionParserConfig) expressionParser {
	return func(request expressionParserRequest) ExpressionParserResult {
		size, err := parseOneOf(request, configs...)
		return ExpressionParserResult{size: size, error: err}
	}
}

func parseOneOf(request expressionParserRequest, configs ...expressionParserConfig) (int, error) {
	if len(configs) <= 0 {
		return 0, nil
	}

	bestResult := ExpressionParserResult{
		error: fmt.Errorf("Unexpected token: %s", request.tokens[0].tokenType),
		size:  -1,
	}

	bestResultIndex := -1

	for i, config := range configs {
		log("parseOneOf i = %d", i)
		result := config.parser(request)

		if result.error == nil {
			config.callback(result)
			return result.size, nil
		}

		if result.size > bestResult.size {
			bestResult = result
			bestResultIndex = i
		}
	}

	configs[bestResultIndex].callback(bestResult)
	return bestResult.size, bestResult.error
}

func functionParser(request expressionParserRequest) ExpressionParserResult {
	function := &Function{}

	size, err := parseAllOrdered(
		request,
		requiredToken(TOKEN_TYPE_KEYWORD_FUNC),
		requiredToken(TOKEN_TYPE_WHITESPACE),
		requiredTokenWithCallback(TOKEN_TYPE_IDENTIFIER,
			func(result ExpressionParserResult) {
				function.name = result.expression.token.value
			}),
		requiredToken(TOKEN_TYPE_ROUND_BRACKET_OPEN),
		optionalToken(TOKEN_TYPE_WHITESPACE),
		requiredToken(TOKEN_TYPE_ROUND_BRACKET_CLOSE),
		optionalToken(TOKEN_TYPE_WHITESPACE),
		requiredToken(TOKEN_TYPE_CURLY_BRACKET_OPEN),
		optionalToken(TOKEN_TYPE_WHITESPACE),
		requiredToken(TOKEN_TYPE_CURLY_BRACKET_CLOSE),
	)

	return ExpressionParserResult{
		size:  size,
		error: err,
		expression: &Expression{
			function: function,
		},
	}
}

func parameterParser(request expressionParserRequest) ExpressionParserResult {
	parameter := &Parameter{}

	size, err := parseAllOrdered(
		request,
		optionalToken(TOKEN_TYPE_WHITESPACE),
		requiredTokenWithCallback(TOKEN_TYPE_IDENTIFIER,
			func(result ExpressionParserResult) {
				parameter.name = result.expression.token.value
			}),
		requiredToken(TOKEN_TYPE_WHITESPACE),
		requiredTokenWithCallback(TOKEN_TYPE_IDENTIFIER,
			func(result ExpressionParserResult) {
				parameter.paramType = result.expression.token.value
			}),
		requiredToken(TOKEN_TYPE_COMMA),
	)

	return ExpressionParserResult{
		size:  size,
		error: err,
		expression: &Expression{
			parameter: parameter,
		},
	}
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
