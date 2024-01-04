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
	const (
		functionNameKey = iota + 1
		functionParamKey
	)

	return allOf(
		requiredToken(token.TYPE_KEYWORD_FUNC),
		requiredToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_IDENTIFIER).emit(functionNameKey),
		requiredToken(token.TYPE_ROUND_BRACKET_OPEN),
		optionalToken(token.TYPE_WHITESPACE),
		parameter().toOptional().toRepeated().emit(functionParamKey),
		requiredToken(token.TYPE_ROUND_BRACKET_CLOSE),
		optionalToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_CURLY_BRACKET_OPEN),
		optionalToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_CURLY_BRACKET_CLOSE),
	).withExpression(func() *Expression { return &Expression{function: &Function{}} }).listen(func(e *Expression, key int, emitted *Expression) {
		switch key {
		case functionNameKey:
			e.function.Name = emitted.token.Value
		case functionParamKey:
			e.function.Parameters = append(e.function.Parameters, emitted.parameter)
		}
	}).withDebug("function")
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
