package ast

import "strings"

type File struct {
	Functions []*Function
}

func file() parser {
	const functionKey = 1
	return oneOf(function().emit(functionKey)).withExpression(func() *Expression { return &Expression{file: &File{}} }).listen(func(e *Expression, key int, emitted *Expression) {
		switch key {
		case functionKey:
			e.file.Functions = append(e.file.Functions, emitted.function)
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
