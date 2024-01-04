package ast

import (
	"github.com/rdidukh/pineapl/token"
)

// TODO: accept required only?
// TODO: disallow ambiguous first token?
func oneOf(parsers ...parser) parser {
	firstTokenTypes := []token.Type{}
	firstTokenMap := map[token.Type]parser{}
	for _, p := range parsers {
		firstTokenTypes = append(firstTokenTypes, p.firstTokenTypes...)
		for _, tokenType := range p.firstTokenTypes {
			firstTokenMap[tokenType] = p
		}
	}

	return parser{
		parserFunc: func(request parserRequest) parserResult {
			it := request.iterator

			token := it.Token()

			p := firstTokenMap[token.Type]

			return p.parse(request)
		},
		firstTokenTypes: firstTokenTypes,
	}
}

func allOf(parsers ...parser) parser {
	firstTokenTypes := []token.Type{}
	for _, p := range parsers {
		firstTokenTypes = append(firstTokenTypes, p.firstTokenTypes...)

		if !p.optional {
			break
		}
	}

	return parser{
		parserFunc: func(request parserRequest) parserResult {
			mergedResult := parserResult{}

			for _, parser := range parsers {
				result := parser.parse(request)

				if result.error != nil {
					return result
				}

				mergedResult.taggedExpressions = append(mergedResult.taggedExpressions, result.taggedExpressions...)
			}

			return mergedResult
		},
		firstTokenTypes: firstTokenTypes,
	}
}
