package ast

import "github.com/rdidukh/pineapl/token"

type CodeBlock struct {
	FunctionCalls []*FunctionCall
}

func codeBlock() parser {
	const (
		functionCallTag = iota + 1
	)
	return allOf(
		requiredToken(token.TYPE_CURLY_BRACKET_OPEN),
		optionalToken(token.TYPE_WHITESPACE),
		functionCall().withTag(functionCallTag).toOptional(),
		optionalToken(token.TYPE_WHITESPACE),
		requiredToken(token.TYPE_CURLY_BRACKET_CLOSE),
	).withExpression(
		func() *Expression { return &Expression{codeBlock: &CodeBlock{}} },
	).listen(
		func(e *Expression, tag int, te *Expression) {
			switch tag {
			case functionCallTag:
				e.codeBlock.FunctionCalls = append(e.codeBlock.FunctionCalls, te.functionCall)
			}
		},
	)
}
