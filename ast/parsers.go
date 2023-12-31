package ast

import (
	"fmt"

	"github.com/rdidukh/pineapl/token"
)

type oneOfParser struct {
	configs []parserConfig
}

func oneOf(configs ...parserConfig) parserConfig {
	parser := oneOfParser{configs: configs}
	return parserConfig{
		parser: parser.parse,
	}
}

func (p oneOfParser) parse(request parserRequest) parserResult {
	size, err := parseOneOf(request, p.configs...)
	return parserResult{
		size:  size,
		error: err,
	}
}

type untilParser struct {
	config     parserConfig
	terminator token.Type
}

func until(config parserConfig, terminator token.Type) parserConfig {
	parser := untilParser{config: config, terminator: terminator}
	return parserConfig{
		parser: parser.parse,
	}
}

func (p untilParser) parse(request parserRequest) parserResult {
	size, err := parseUntil(request, p.terminator, p.config)
	return parserResult{
		size:  size,
		error: err,
	}
}

type allOfParser struct {
	configs []parserConfig
}

func allOf(configs ...parserConfig) parserConfig {
	parser := allOfParser{configs: configs}
	return parserConfig{
		parser: parser.parse,
	}
}

func (p allOfParser) parse(request parserRequest) parserResult {
	size, err := parseAllOrdered(request, p.configs...)
	return parserResult{
		size:  size,
		error: err,
	}
}

// TODO: inline.
func parseOneOf(request parserRequest, configs ...parserConfig) (int, error) {
	if len(configs) <= 0 {
		return 0, nil
	}

	bestResult := parserResult{
		error: fmt.Errorf("unexpected token: %s", request.tokens[0].Type),
		size:  -1,
	}

	for _, config := range configs {
		result := config.parser(request)

		if result.error == nil {
			return result.size, nil
		}

		if result.size > bestResult.size {
			bestResult = result
		}
	}

	return bestResult.size, bestResult.error
}

// TODO: inline.
func parseUntil(request parserRequest, terminator token.Type, config parserConfig) (int, error) {
	offset := 0
	tokens := request.tokens

	for offset < len(tokens) && tokens[offset].Type != terminator {
		size, err := parseOneOf(parserRequest{
			tokens: request.tokens[offset:],
		}, config)

		if err != nil {
			return size, err
		}

		offset += size
	}

	if offset >= len(tokens) {
		return offset, fmt.Errorf("expected %s, found: EOF", terminator)
	}

	nextTokenType := tokens[offset].Type

	if nextTokenType != terminator {
		return offset, fmt.Errorf("expected %s, found: %s", terminator, nextTokenType)
	}

	return offset + 1, nil
}

// TODO: inline.
func parseAllOrdered(request parserRequest, configs ...parserConfig) (int, error) {
	offset := 0
	for _, config := range configs {
		result := config.parser(parserRequest{
			tokens: request.tokens[offset:],
		})

		offset += result.size

		if result.error != nil {
			return offset, result.error
		}
	}

	return offset, nil
}
