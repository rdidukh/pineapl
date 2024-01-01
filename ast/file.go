package ast

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
