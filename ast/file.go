package ast

type File struct {
	Functions []*Function
}

func file() parserConfig {
	file := &File{}
	return oneOf(function().withCallback(
		func(r parserResult) {
			file.Functions = append(file.Functions, r.expression.function)
		},
	)).withExpression(&Expression{file: file})
}
