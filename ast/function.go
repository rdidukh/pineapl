package ast

import (
	"github.com/rdidukh/pineapl/token"
)

type Function struct {
	Name       string
	Parameters []*Parameter
}

func functionParser(request parserRequest) parserResult {
	function := &Function{}

	size, err := parseAllOrdered(
		request,
		requiredToken(token.TYPE_KEYWORD_FUNC),
		requiredToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_IDENTIFIER).withCallback(
			func(result parserResult) {
				function.Name = result.expression.token.Value
			}),
		requiredToken(token.TYPE_ROUND_BRACKET_OPEN),
		optionalToken(token.TYPE_WHITESPACE),
		parserConfig{
			parser: oneOfRepeatedUntilParser(
				token.TYPE_ROUND_BRACKET_CLOSE,
				parserConfig{
					parser: parameterParser,
					callback: func(result parserResult) {
						function.Parameters = append(function.Parameters, result.expression.parameter)
					},
				},
			),
		},
		optionalToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_CURLY_BRACKET_OPEN),
		optionalToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_CURLY_BRACKET_CLOSE),
	)

	return parserResult{
		size:  size,
		error: err,
		expression: &Expression{
			function: function,
		},
	}
}
