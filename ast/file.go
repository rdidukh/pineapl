package ast

type File struct {
	Functions []*Function
}

func fileParser(request parserRequest) parserResult {
	file := &File{}

	size, err := parseOneOfRepeated(request,
		parserConfig{
			parser: functionParser,
			callback: func(result parserResult) {
				file.Functions = append(file.Functions, result.expression.function)
			},
		})

	return parserResult{size: size, error: err, expression: &Expression{file: file}}
}
