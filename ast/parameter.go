package ast

import (
	"github.com/rdidukh/pineapl/token"
)

type Parameter struct {
	Name string
	Type string
}

func parameter() parserConfig {
	// TODO: this is incorrect as this can be used multiple times.
	param := &Parameter{}
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
	).withExpression(&Expression{parameter: param})
}
