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

type parserFunc func(parserRequest) parserResult
type parserCallback func(result parserResult)

type parser struct {
	parserFunc parserFunc
}

func (p parser) withCallback(callback parserCallback) parser {
	parser := p
	parser.parserFunc = func(request parserRequest) parserResult {
		result := p.parserFunc(request)
		if result.error == nil {
			callback(result)
		}
		return result
	}
	return parser
}

// TODO: find a thread safe (stateless) way of reusing a parser.
func (p parser) withInit(setUpFunc func()) parser {
	parser := p
	parser.parserFunc = func(request parserRequest) parserResult {
		setUpFunc()
		return p.parserFunc(request)
	}
	return parser
}

func (p parser) withDebug(debug string) parser {
	parser := p
	parser.parserFunc = func(request parserRequest) parserResult {
		logger.LogPadded(debugPadding, "Before calling parser %s %d", debug, len(request.tokens))
		debugPadding += 1
		result := p.parserFunc(request)
		debugPadding -= 1
		logger.LogPadded(debugPadding, "After calling parser %s expr=%v", debug, result.expression)
		return result
	}
	return parser
}

func (p parser) withExpression(exprFunc func() *Expression) parser {
	parser := p
	parser.parserFunc = func(request parserRequest) parserResult {
		result := p.parserFunc(request)
		result.expression = exprFunc()
		return result
	}
	return parser
}

func ParseString(code string) ([]*token.Token, *File, error) {
	tokens, err := token.GetTokens(code)

	if err != nil {
		return tokens, nil, err
	}

	result := file().parserFunc(parserRequest{tokens: tokens})

	return tokens, result.expression.file, result.error
}

func Codegen(code string) (string, error) {
	_, file, err := ParseString(code)
	if err != nil {
		return "", err
	}

	return file.Codegen()
}
