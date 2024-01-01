package ast

import (
	"github.com/rdidukh/pineapl/token"
)

type Parameter struct {
	Name string
	Type string
}

func parameter() parser {
	var param *Parameter
	return allOf(
		optionalToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_IDENTIFIER).withCallback(
			func(r parserResult) {
				param.Name = r.expression.token.Value
			}),
		requiredToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_IDENTIFIER).withCallback(
			func(r parserResult) {
				param.Type = r.expression.token.Value
			}),
		requiredToken(token.TYPE_COMMA),
	).withInit(func() {
		param = &Parameter{}
	}).withExpression(func() *Expression { return &Expression{parameter: param} })
}
