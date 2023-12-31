package ast

import (
	"github.com/rdidukh/pineapl/logger"
	"github.com/rdidukh/pineapl/token"
)

var debugPadding = 0

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
	parser parser
}

func (c parserConfig) withCallback(callback parserCallback) parserConfig {
	config := c
	config.parser = func(request parserRequest) parserResult {
		result := c.parser(request)
		if result.error == nil {
			callback(result)
		}
		return result
	}
	return config
}

func (c parserConfig) withDebug(debug string) parserConfig {
	config := c
	config.parser = func(request parserRequest) parserResult {
		logger.LogPadded(debugPadding, "Before calling parser %s %d", debug, len(request.tokens))
		debugPadding += 1
		result := c.parser(request)
		debugPadding -= 1
		logger.LogPadded(debugPadding, "After calling parser %s expr=%v", debug, result.expression)
		return result
	}
	return config
}

func (c parserConfig) withExpression(expression *Expression) parserConfig {
	config := c
	config.parser = func(request parserRequest) parserResult {
		result := c.parser(request)
		result.expression = expression
		return result
	}
	return config
}

func ParseString(code string) ([]*token.Token, *File, error) {
	tokens, err := token.GetTokens(code)

	if err != nil {
		return tokens, nil, err
	}

	result := file().parser(parserRequest{tokens: tokens})

	return tokens, result.expression.file, result.error
}
