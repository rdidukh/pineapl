package ast

import (
	"fmt"

	"github.com/rdidukh/pineapl/token"
)

type requiredTokenParser struct {
	tokenType token.Type
}

func requiredToken(tokenType token.Type) parserConfig {
	parser := requiredTokenParser{tokenType: tokenType}
	return parserConfig{
		parser: parser.parse,
	}
}

func (p requiredTokenParser) parse(request parserRequest) parserResult {
	tokens := request.tokens
	if len(tokens) <= 0 {
		return parserResult{
			error: fmt.Errorf("expected %s, found: EOF", p.tokenType),
		}
	}

	actualTokenType := tokens[0].Type

	if actualTokenType != p.tokenType {
		return parserResult{
			error: fmt.Errorf("expected %s, found: %s", p.tokenType, actualTokenType),
		}
	}

	return parserResult{
		expression: &Expression{token: tokens[0]},
		size:       1,
	}
}

type optionalTokenParser struct {
	tokenType token.Type
}

func optionalToken(tokenType token.Type) parserConfig {
	parser := optionalTokenParser{tokenType: tokenType}
	return parserConfig{
		parser: parser.parse,
	}
}

func (p optionalTokenParser) parse(request parserRequest) parserResult {
	tokens := request.tokens
	if len(tokens) <= 0 || tokens[0].Type != p.tokenType {
		return parserResult{}
	}

	return parserResult{
		expression: &Expression{token: tokens[0]},
		size:       1,
	}
}
