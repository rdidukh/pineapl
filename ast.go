package main

import (
	"fmt"

	"github.com/rdidukh/pineapl/token"
)

type Expression struct {
	token     *token.Token
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

type expressionParserRequest struct {
	tokens   []*token.Token
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

func ParseFile(tokens []*token.Token) (*File, error) {
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

func requiredToken(tokenType token.Type) expressionParserConfig {
	return requiredTokenWithCallback(tokenType, func(result ExpressionParserResult) {})
}

func requiredTokenWithCallback(tokenType token.Type, callback expressionParserCallback) expressionParserConfig {
	return expressionParserConfig{
		parser:   requiredTokenParser(tokenType),
		callback: callback,
	}
}

func requiredTokenParser(tokenType token.Type) expressionParser {
	return func(request expressionParserRequest) ExpressionParserResult {
		tokens := request.tokens
		if len(tokens) <= 0 {
			return ExpressionParserResult{
				error: fmt.Errorf("Expected %s, Found: EOF", tokenType),
			}
		}

		actualTokenType := tokens[0].Type

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

func optionalToken(tokenType token.Type) expressionParserConfig {
	return expressionParserConfig{
		parser: optionalTokenParser(tokenType),
	}
}

func optionalTokenParser(tokenType token.Type) expressionParser {
	return func(request expressionParserRequest) ExpressionParserResult {
		tokens := request.tokens
		if len(tokens) <= 0 || tokens[0].Type != tokenType {
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
	return parseOneOfRepeatedUntil(request, token.TYPE_EOF, configs...)
}

func oneOfRepeatedUntilParser(terminator token.Type, configs ...expressionParserConfig) expressionParser {
	return func(request expressionParserRequest) ExpressionParserResult {
		size, err := parseOneOfRepeatedUntil(request, terminator, configs...)
		return ExpressionParserResult{
			size:  size,
			error: err,
		}
	}
}

func parseOneOfRepeatedUntil(request expressionParserRequest, terminator token.Type, configs ...expressionParserConfig) (int, error) {
	offset := 0
	tokens := request.tokens

	for offset < len(tokens) && tokens[offset].Type != terminator {
		size, err := parseOneOf(expressionParserRequest{
			tokens: request.tokens[offset:],
		}, configs...)

		if err != nil {
			return size, err
		}

		offset += size
	}

	if offset >= len(tokens) {
		return offset, fmt.Errorf("expected %s, found: EOF", terminator)
	}

	nextTokenType := tokens[offset].Type

	if nextTokenType != terminator {
		return offset, fmt.Errorf("expected %s, found: %s", terminator, nextTokenType)
	}

	return offset + 1, nil
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
		error: fmt.Errorf("Unexpected token: %s", request.tokens[0].Type),
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
		requiredToken(token.TYPE_KEYWORD_FUNC),
		requiredToken(token.TYPE_WHITESPACE),
		requiredTokenWithCallback(token.TYPE_IDENTIFIER,
			func(result ExpressionParserResult) {
				function.name = result.expression.token.Value
			}),
		requiredToken(token.TYPE_ROUND_BRACKET_OPEN),
		optionalToken(token.TYPE_WHITESPACE),
		expressionParserConfig{
			parser: oneOfRepeatedUntilParser(
				token.TYPE_ROUND_BRACKET_CLOSE,
				expressionParserConfig{
					parser: parameterParser,
					callback: func(result ExpressionParserResult) {
						function.parameters = append(function.parameters, result.expression.parameter)
					},
				},
			),
		},
		optionalToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_CURLY_BRACKET_OPEN),
		optionalToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_CURLY_BRACKET_CLOSE),
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
		optionalToken(token.TYPE_WHITESPACE),
		requiredTokenWithCallback(token.TYPE_IDENTIFIER,
			func(result ExpressionParserResult) {
				parameter.name = result.expression.token.Value
			}),
		requiredToken(token.TYPE_WHITESPACE),
		requiredTokenWithCallback(token.TYPE_IDENTIFIER,
			func(result ExpressionParserResult) {
				parameter.paramType = result.expression.token.Value
			}),
		requiredToken(token.TYPE_COMMA),
	)

	return ExpressionParserResult{
		size:  size,
		error: err,
		expression: &Expression{
			parameter: parameter,
		},
	}
}
