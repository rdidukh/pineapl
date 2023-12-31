package ast

import (
	"github.com/rdidukh/pineapl/token"
)

type Function struct {
	Name       string
	Parameters []*Parameter
}

func function() parserConfig {
	function := &Function{}
	return allOf(
		requiredToken(token.TYPE_KEYWORD_FUNC),
		requiredToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_IDENTIFIER).withCallback(
			func(r parserResult) {
				function.Name = r.expression.token.Value
			}),
		requiredToken(token.TYPE_ROUND_BRACKET_OPEN),
		optionalToken(token.TYPE_WHITESPACE),
		until(parameter().withCallback(
			func(r parserResult) {
				function.Parameters = append(function.Parameters, r.expression.parameter)
			},
		), token.TYPE_ROUND_BRACKET_CLOSE),
		optionalToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_CURLY_BRACKET_OPEN),
		optionalToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_CURLY_BRACKET_CLOSE),
	).withExpression(&Expression{function: function})
}
