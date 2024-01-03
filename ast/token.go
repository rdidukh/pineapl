package ast

import (
	"fmt"

	"github.com/rdidukh/pineapl/token"
)

func requiredToken(tokenType token.Type) parser {
	return parser{
		parserFunc: func(request parserRequest) parserResult {
			return parseToken(request, tokenType)
		},
		firstTokenTypes: []token.Type{tokenType},
		debug:           fmt.Sprintf("token %s", tokenType),
	}
}

func optionalToken(tokenType token.Type) parser {
	return requiredToken(tokenType).toOptional()
}

func parseToken(request parserRequest, tokenType token.Type) parserResult {
	it := request.iterator
	token := it.Token()

	if token.Type != tokenType {
		return parserResult{
			error: fmt.Errorf("expected %s, found: %s(%q)", tokenType, token.Type, token.Value),
		}
	}

	it.Advance()

	return parserResult{
		expression: &Expression{token: token},
	}
}
