package ast

import (
	"fmt"

	"github.com/rdidukh/pineapl/token"
)

func oneOf(parsers ...parser) parser {
	return parser{
		parserFunc: func(request parserRequest) parserResult {
			return parseOneOf(request, parsers...)
		},
	}
}

func until(p parser, terminator token.Type) parser {
	return parser{
		parserFunc: func(request parserRequest) parserResult {
			return parseUntil(request, terminator, p)
		},
	}
}

func allOf(parsers ...parser) parser {
	return parser{
		parserFunc: func(request parserRequest) parserResult {
			size, err := parseAllOrdered(request, parsers...)
			return parserResult{
				size:  size,
				error: err,
			}
		},
	}
}

// TODO: inline.
func parseOneOf(request parserRequest, parsers ...parser) parserResult {
	if len(parsers) <= 0 {
		return parserResult{}
	}

	bestResult := parserResult{
		error: fmt.Errorf("unexpected token: %s", request.tokens[0].Type),
		size:  -1,
	}

	for _, parser := range parsers {
		result := parser.parserFunc(request)

		if result.error == nil {
			return result
		}

		if result.size > bestResult.size {
			bestResult = result
		}
	}

	return bestResult
}

// TODO: inline.
func parseUntil(request parserRequest, terminator token.Type, p parser) parserResult {
	offset := 0
	tokens := request.tokens

	for offset < len(tokens) && tokens[offset].Type != terminator {
		result := parseOneOf(parserRequest{
			tokens: request.tokens[offset:],
		}, p)

		if result.error != nil {
			return result
		}

		offset += result.size
	}

	if offset >= len(tokens) {
		return parserResult{
			size:  offset,
			error: fmt.Errorf("expected %s, found: EOF", terminator),
		}
	}

	nextTokenType := tokens[offset].Type

	if nextTokenType != terminator {
		return parserResult{
			size:  offset,
			error: fmt.Errorf("expected %s, found: %s", terminator, nextTokenType),
		}
	}

	return parserResult{size: offset + 1}
}

// TODO: inline.
func parseAllOrdered(request parserRequest, parsers ...parser) (int, error) {
	offset := 0
	for _, parser := range parsers {
		result := parser.parserFunc(parserRequest{
			tokens: request.tokens[offset:],
		})

		offset += result.size

		if result.error != nil {
			return offset, result.error
		}
	}

	return offset, nil
}
