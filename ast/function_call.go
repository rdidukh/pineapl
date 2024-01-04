package ast

import (
	"github.com/rdidukh/pineapl/token"
)

type FunctionCall struct {
	Name string
}

func functionCall() parser {
	const (
		functionNameTag = iota + 1
	)

	return allOf(
		requiredToken(token.TYPE_IDENTIFIER).withTag(functionNameTag),
		requiredToken(token.TYPE_ROUND_BRACKET_OPEN),
		requiredToken(token.TYPE_ROUND_BRACKET_CLOSE),
	).withExpression(
		func() *Expression { return &Expression{functionCall: &FunctionCall{}} },
	).listen(
		func(e *Expression, tag int, te *Expression) {
			switch tag {
			case functionNameTag:
				e.functionCall.Name = te.token.Value
			}
		},
	)
}
