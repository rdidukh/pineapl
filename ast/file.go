package ast

import "strings"

type File struct {
	Functions []*Function
}

func file() parser {
	const functionTag = 1
	return oneOf(
		function().withTag(functionTag),
	).withExpression(
		func() *Expression {
			return &Expression{file: &File{}}
		},
	).listen(
		func(e *Expression, tag int, te *Expression) {
			switch tag {
			case functionTag:
				e.file.Functions = append(e.file.Functions, te.function)
			}
		}).withDebug("file")
}

func (f *File) Codegen() (string, error) {
	code := strings.Builder{}
	for _, f := range f.Functions {
		code.WriteString(f.codegen())
	}
	return code.String(), nil
}
