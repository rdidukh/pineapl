package ast

import (
	"fmt"
	"strings"

	"github.com/rdidukh/pineapl/token"
)

type Function struct {
	Name       string
	Parameters []*Parameter
	CodeBlock  *CodeBlock
}

func function() parser {
	const (
		functionNameTag = iota + 1
		functionParamTag
		functionCodeBlockTag
	)

	return allOf(
		requiredToken(token.TYPE_KEYWORD_FUNC),
		requiredToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_IDENTIFIER).withTag(functionNameTag),
		requiredToken(token.TYPE_ROUND_BRACKET_OPEN),
		optionalToken(token.TYPE_WHITESPACE),
		parameter().toOptional().toRepeated().withTag(functionParamTag),
		requiredToken(token.TYPE_ROUND_BRACKET_CLOSE),
		optionalToken(token.TYPE_WHITESPACE),
		codeBlock().withTag(functionCodeBlockTag),
	).withExpression(
		func() *Expression { return &Expression{function: &Function{}} },
	).listen(
		func(e *Expression, tag int, te *Expression) {
			switch tag {
			case functionNameTag:
				e.function.Name = te.token.Value
			case functionParamTag:
				e.function.Parameters = append(e.function.Parameters, te.parameter)
			case functionCodeBlockTag:
				e.function.CodeBlock = te.codeBlock
			}
		},
	).withDebug("function")
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
