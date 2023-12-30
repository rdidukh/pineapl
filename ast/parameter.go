package ast

import "github.com/rdidukh/pineapl/token"

type Parameter struct {
	Name string
	Type string
}

func parameterParser(request parserRequest) parserResult {
	parameter := &Parameter{}

	size, err := parseAllOrdered(
		request,
		optionalToken(token.TYPE_WHITESPACE),
		requiredTokenWithCallback(token.TYPE_IDENTIFIER,
			func(result parserResult) {
				parameter.Name = result.expression.token.Value
			}),
		requiredToken(token.TYPE_WHITESPACE),
		requiredTokenWithCallback(token.TYPE_IDENTIFIER,
			func(result parserResult) {
				parameter.Type = result.expression.token.Value
			}),
		requiredToken(token.TYPE_COMMA),
	)

	return parserResult{
		size:  size,
		error: err,
		expression: &Expression{
			parameter: parameter,
		},
	}
}
