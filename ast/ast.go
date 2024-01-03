package ast

import (
	"fmt"
	"slices"

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
	iterator *token.Iterator
}

type parserResult struct {
	error      error
	expression *Expression
}

type parserFunc func(parserRequest) parserResult
type parserCallback func(result parserResult)

type parser struct {
	parserFunc      parserFunc
	firstTokenTypes []token.Type
	optional        bool
	repeated        bool
	initFunc        func()
	callback        parserCallback
	debug           string
	expressionFunc  func() *Expression
}

func (p parser) toOptional() parser {
	result := p
	result.optional = true
	return result
}

func (p parser) toRepeated() parser {
	result := p
	result.repeated = true
	return result
}

func (p parser) parse(request parserRequest) parserResult {
	it := request.iterator
	token := it.Token()
	eof := it.IsEof()
	matchesToken := slices.Contains(p.firstTokenTypes, token.Type) // TODO && !eof
	if !p.optional && (!matchesToken || eof) {
		return parserResult{
			error: fmt.Errorf("unexpected token %s(%q), expected: %s", token.Type, token.Value, p.firstTokenTypes),
		}
	}

	if p.optional && (!matchesToken || eof) {
		return parserResult{}
	}

	if !matchesToken {
		panic("!matchesToken")
	}

	if p.initFunc != nil {
		p.initFunc()
	}

	if p.debug != "" {
		logger.LogPadded(debugPadding, "Before calling parser %s", p.debug)
		debugPadding += 1
	}
	result := p.parserFunc(request)
	if p.debug != "" {
		logger.LogPadded(debugPadding, "After calling parser %s expr=%v", p.debug, result.expression)
		debugPadding -= 1
	}

	if result.error != nil {
		return result
	}

	if p.expressionFunc != nil {
		result.expression = p.expressionFunc()
	}

	if p.callback != nil {
		p.callback(result)
	}

	if !p.repeated {
		return result
	}

	// TODO: optimize recursion.
	return p.toOptional().parse(request)
}

func (p parser) withCallback(callback parserCallback) parser {
	parser := p
	parser.callback = callback
	return parser
}

// TODO: find a thread safe (stateless) way of reusing a parser.
func (p parser) withInit(setUpFunc func()) parser {
	parser := p
	parser.initFunc = setUpFunc
	return parser
}

func (p parser) withDebug(debug string) parser {
	parser := p
	p.debug = debug
	return parser
}

func (p parser) withExpression(exprFunc func() *Expression) parser {
	parser := p
	parser.expressionFunc = exprFunc
	return parser
}

func ParseString(code string) (*File, error) {
	iterator, err := token.GetIterator(code)

	if err != nil {
		return nil, err
	}

	result := file().parse(parserRequest{iterator: iterator})

	return result.expression.file, result.error
}

func Codegen(code string) (string, error) {
	file, err := ParseString(code)
	if err != nil {
		return "", err
	}

	return file.Codegen()
}
