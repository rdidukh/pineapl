package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
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

func (f *Function) addToModule(m *ir.Module) {
	params := []*ir.Param{}

	for _, param := range f.Parameters {
		typ := Type(param.Type).toIrType()
		// TOOD: move to parameters.go
		params = append(params, ir.NewParam(param.Name, typ))
	}

	retType := types.I1

	result := m.NewFunc(f.Name, retType, params...)

	result.NewBlock("").NewRet(constant.NewInt(retType, 0))
}
