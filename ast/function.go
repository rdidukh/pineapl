package ast

import (
	"fmt"
	"strings"

	"github.com/rdidukh/pineapl/token"
)

type Function struct {
	Name       string
	Parameters []*Parameter
}

func function() parser {
	var function *Function
	return allOf(
		requiredToken(token.TYPE_KEYWORD_FUNC),
		requiredToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_IDENTIFIER).withCallback(
			func(r parserResult) {
				function.Name = r.expression.token.Value
			}),
		requiredToken(token.TYPE_ROUND_BRACKET_OPEN),
		optionalToken(token.TYPE_WHITESPACE),
		parameter().toOptional().toRepeated().withCallback(func(r parserResult) {
			function.Parameters = append(function.Parameters, r.expression.parameter)
		}),
		requiredToken(token.TYPE_ROUND_BRACKET_CLOSE),
		optionalToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_CURLY_BRACKET_OPEN),
		optionalToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_CURLY_BRACKET_CLOSE),
	).withInit(func() {
		function = &Function{}
	}).withExpression(func() *Expression { return &Expression{function: function} })
}

func (f *Function) codegen() string {
	code := strings.Builder{}

	code.WriteString(fmt.Sprintf("define void @%s(", f.Name))
	for i, param := range f.Parameters {
		code.WriteString(fmt.Sprintf("%s %s", param.Type, param.Name))
		if i+1 < len(f.Parameters) {
			code.WriteString(", ")
		}
	}
	code.WriteString(") {\n}\n")

	return code.String()
}
