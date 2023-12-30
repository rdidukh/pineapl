package ast

import (
	"fmt"

	"github.com/rdidukh/pineapl/token"
)

func requiredToken(tokenType token.Type) parserConfig {
	return parserConfig{
		parser: requiredTokenParser(tokenType),
	}
}

func requiredTokenParser(tokenType token.Type) parser {
	return func(request parserRequest) parserResult {
		tokens := request.tokens
		if len(tokens) <= 0 {
			return parserResult{
				error: fmt.Errorf("expected %s, found: EOF", tokenType),
			}
		}

		actualTokenType := tokens[0].Type

		if actualTokenType != tokenType {
			return parserResult{
				error: fmt.Errorf("expected %s, found: %s", tokenType, actualTokenType),
			}
		}

		return parserResult{
			expression: &Expression{token: tokens[0]},
			size:       1,
		}
	}
}

func optionalToken(tokenType token.Type) parserConfig {
	return parserConfig{
		parser: optionalTokenParser(tokenType),
	}
}

func optionalTokenParser(tokenType token.Type) parser {
	return func(request parserRequest) parserResult {
		tokens := request.tokens
		if len(tokens) <= 0 || tokens[0].Type != tokenType {
			return parserResult{}
		}

		return parserResult{
			expression: &Expression{token: tokens[0]},
			size:       1,
		}
	}
}
