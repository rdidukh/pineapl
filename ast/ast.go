package ast

import (
	"fmt"

	"github.com/rdidukh/pineapl/logger"
	"github.com/rdidukh/pineapl/token"
)

type Expression struct {
	token     *token.Token
	function  *Function
	file      *File
	parameter *Parameter
}

type parserRequest struct {
	tokens []*token.Token
}

type parserResult struct {
	size       int
	error      error
	expression *Expression
}

type parser func(parserRequest) parserResult
type parserCallback func(result parserResult)

type parserConfig struct {
	parser   parser
	callback parserCallback
}

func (c parserConfig) onSuccess(result parserResult) {
	if c.callback != nil {
		c.callback(result)
	}
}

func ParseString(code string) ([]*token.Token, *File, error) {
	tokens, err := token.GetTokens(code)

	if err != nil {
		return tokens, nil, err
	}

	file, err := ParseFile(tokens)

	return tokens, file, err

}

func ParseFile(tokens []*token.Token) (*File, error) {
	result := fileParser(parserRequest{tokens: tokens})

	return result.expression.file, result.error
}

func requiredToken(tokenType token.Type) parserConfig {
	return requiredTokenWithCallback(tokenType, func(result parserResult) {})
}

func requiredTokenWithCallback(tokenType token.Type, callback parserCallback) parserConfig {
	return parserConfig{
		parser:   requiredTokenParser(tokenType),
		callback: callback,
	}
}

func requiredTokenParser(tokenType token.Type) parser {
	return func(request parserRequest) parserResult {
		tokens := request.tokens
		if len(tokens) <= 0 {
			return parserResult{
				error: fmt.Errorf("Expected %s, Found: EOF", tokenType),
			}
		}

		actualTokenType := tokens[0].Type

		if actualTokenType != tokenType {
			return parserResult{
				error: fmt.Errorf("Expected %s, Found: %s", tokenType, actualTokenType),
			}
		}

		return parserResult{
			expression: &Expression{token: tokens[0]},
			size:       1,
		}
	}
}

func optionalToken(tokenType token.Type) parserConfig {
	return parserConfig{
		parser: optionalTokenParser(tokenType),
	}
}

func optionalTokenParser(tokenType token.Type) parser {
	return func(request parserRequest) parserResult {
		tokens := request.tokens
		if len(tokens) <= 0 || tokens[0].Type != tokenType {
			return parserResult{}
		}

		return parserResult{
			expression: &Expression{token: tokens[0]},
			size:       1,
		}
	}
}

func oneOfRepeatedParser(configs ...parserConfig) parser {
	return func(request parserRequest) parserResult {
		size, err := parseOneOfRepeated(request, configs...)
		return parserResult{
			size:  size,
			error: err,
		}
	}
}

func parseOneOfRepeated(request parserRequest, configs ...parserConfig) (int, error) {
	return parseOneOfRepeatedUntil(request, token.TYPE_EOF, configs...)
}

func oneOfRepeatedUntilParser(terminator token.Type, configs ...parserConfig) parser {
	return func(request parserRequest) parserResult {
		size, err := parseOneOfRepeatedUntil(request, terminator, configs...)
		return parserResult{
			size:  size,
			error: err,
		}
	}
}

func parseOneOfRepeatedUntil(request parserRequest, terminator token.Type, configs ...parserConfig) (int, error) {
	offset := 0
	tokens := request.tokens

	for offset < len(tokens) && tokens[offset].Type != terminator {
		size, err := parseOneOf(parserRequest{
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

func allOrderedParser(configs ...parserConfig) parser {
	return func(request parserRequest) parserResult {
		size, err := parseAllOrdered(request, configs...)
		return parserResult{
			size:  size,
			error: err,
		}
	}
}

func parseAllOrdered(request parserRequest, configs ...parserConfig) (int, error) {
	offset := 0
	for _, config := range configs {
		result := config.parser(parserRequest{
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

func oneOfParser(configs ...parserConfig) parser {
	return func(request parserRequest) parserResult {
		size, err := parseOneOf(request, configs...)
		return parserResult{size: size, error: err}
	}
}

func parseOneOf(request parserRequest, configs ...parserConfig) (int, error) {
	if len(configs) <= 0 {
		return 0, nil
	}

	bestResult := parserResult{
		error: fmt.Errorf("Unexpected token: %s", request.tokens[0].Type),
		size:  -1,
	}

	bestResultIndex := -1

	for i, config := range configs {
		logger.Log("parseOneOf i = %d", i)
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
