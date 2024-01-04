package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/rdidukh/pineapl/token"
)

type Parameter struct {
	Name string
	Type string
}

func parameter() parser {
	const (
		paramNameTag = iota + 1
		paramTypeTag
	)
	return allOf(
		optionalToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_IDENTIFIER).withTag(paramNameTag),
		requiredToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_IDENTIFIER).withTag(paramTypeTag),
		requiredToken(token.TYPE_COMMA),
	).withExpression(
		func() *Expression { return &Expression{parameter: &Parameter{}} },
	).listen(
		func(e *Expression, tag int, te *Expression) {
			switch tag {
			case paramNameTag:
				e.parameter.Name = te.token.Value
			case paramTypeTag:
				e.parameter.Type = te.token.Value
			}
		},
	).withDebug("parameter")
}

func (p *Parameter) toLlvmParam() *ir.Param {
	typ := Type(p.Type).toIrType()
	return ir.NewParam(p.Name, typ)
}
