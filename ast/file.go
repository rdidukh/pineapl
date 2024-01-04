package ast

import (
	"github.com/llir/llvm/ir"
)

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

func (f *File) addToModule(m *ir.Module) {
	for _, fun := range f.Functions {
		fun.addToModule(m)
	}
}

func (f *File) Codegen() (string, error) {
	module := ir.NewModule()

	f.addToModule(module)

	return module.String(), nil
}
