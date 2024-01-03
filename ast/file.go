package ast

import "strings"

type File struct {
	Functions []*Function
}

func file() parser {
	var file *File
	return oneOf(function().withCallback(
		func(r parserResult) {
			file.Functions = append(file.Functions, r.expression.function)
		},
	)).withInit(func() {
		file = &File{}
	}).withExpression(func() *Expression { return &Expression{file: file} })
}

func (f *File) Codegen() (string, error) {
	code := strings.Builder{}
	for _, f := range f.Functions {
		code.WriteString(f.codegen())
	}
	return code.String(), nil
}
