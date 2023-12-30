package ast

import (
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
