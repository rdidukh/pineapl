package ast

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

func (c parserConfig) withCallback(callback parserCallback) parserConfig {
	config := c
	config.callback = callback
	return config
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
