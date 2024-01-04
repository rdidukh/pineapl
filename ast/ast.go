package ast

import (
	"fmt"
	"slices"

	"github.com/rdidukh/pineapl/logger"
	"github.com/rdidukh/pineapl/token"
)

var debugPadding = 0

type Expression struct {
	token        *token.Token
	function     *Function
	file         *File
	parameter    *Parameter
	codeBlock    *CodeBlock
	functionCall *FunctionCall
}

type parserRequest struct {
	iterator *token.Iterator
}

type taggedExpression struct {
	expression *Expression
	tag        int
}

type parserResult struct {
	// Set by parseFunc
	error      error
	expression *Expression

	// Set by parser.parse().
	taggedExpressions []taggedExpression
}

type parserFunc func(parserRequest) parserResult
type listenerFunc func(e *Expression, tag int, taggedExpr *Expression)

type parser struct {
	parserFunc      parserFunc
	firstTokenTypes []token.Type
	optional        bool
	repeated        bool
	debug           string
	tag             int
	exprFunc        func() *Expression
	listenerFunc    listenerFunc
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

func logDebug(tag string, format string, a ...any) {
	if tag == "" {
		return
	}
	logger.LogPadded(debugPadding, fmt.Sprintf("[%s] %s", tag, format), a...)
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

	logDebug(p.debug, "Before calling parserFunc")
	debugPadding++
	result := p.parserFunc(request)
	debugPadding--
	logDebug(p.debug, "After calling parserFunc")

	if result.error != nil {
		return result
	}

	repeatedResult := parserResult{}
	if p.repeated {
		// TODO: optimize recursion.
		repeatedResult = p.toOptional().parse(request)
	}

	if repeatedResult.error != nil {
		return repeatedResult
	}

	expression := result.expression
	if result.expression == nil && p.exprFunc != nil {
		expression = p.exprFunc()

		for _, te := range result.taggedExpressions {
			p.listenerFunc(expression, te.tag, te.expression)
		}
	}

	taggedExpressions := repeatedResult.taggedExpressions

	if p.tag != 0 && expression != nil {
		taggedExpressions = append([]taggedExpression{{expression: expression, tag: p.tag}}, taggedExpressions...)
	}

	finalResult := parserResult{
		taggedExpressions: taggedExpressions,
	}

	return finalResult
}

func (p parser) withTag(tag int) parser {
	parser := p
	parser.tag = tag
	return parser
}

func (p parser) listen(listener listenerFunc) parser {
	parser := p
	parser.listenerFunc = listener
	return parser
}

func (p parser) withExpression(exprFunc func() *Expression) parser {
	parser := p
	parser.exprFunc = exprFunc
	return parser
}

func (p parser) withDebug(debug string) parser {
	parser := p
	parser.debug = debug
	return parser
}

func ParseString(code string) (*File, error) {
	iterator, err := token.GetIterator(code)

	if err != nil {
		return nil, err
	}

	const fileKey = 1

	result := file().withTag(fileKey).parse(parserRequest{iterator: iterator})

	if result.error != nil {
		return nil, result.error
	}

	return result.taggedExpressions[0].expression.file, nil
}

func Codegen(code string) (string, error) {
	file, err := ParseString(code)
	if err != nil {
		return "", err
	}

	return file.Codegen()
}
