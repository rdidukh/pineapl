package ast

import (
	"github.com/rdidukh/pineapl/token"
)

type Parameter struct {
	Name string
	Type string
}

func parameter() parser {
	const (
		paramNameKey = 1
		paramTypeKey = 2
	)
	return allOf(
		optionalToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_IDENTIFIER).emit(paramNameKey),
		requiredToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_IDENTIFIER).emit(paramTypeKey),
		requiredToken(token.TYPE_COMMA),
	).withExpression(func() *Expression { return &Expression{parameter: &Parameter{}} }).listen(
		func(e *Expression, key int, emitted *Expression) {
			switch key {
			case paramNameKey:
				e.parameter.Name = emitted.token.Value
			case paramTypeKey:
				e.parameter.Type = emitted.token.Value
			}
		},
	).withDebug("parameter")
}
