package ast

import (
	"fmt"

	"github.com/rdidukh/pineapl/logger"
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

func parseOneOf(request parserRequest, configs ...parserConfig) (int, error) {
	if len(configs) <= 0 {
		return 0, nil
	}

	bestResult := parserResult{
		error: fmt.Errorf("unexpected token: %s", request.tokens[0].Type),
		size:  -1,
	}

	bestResultIndex := -1

	for i, config := range configs {
		logger.Log("parseOneOf i = %d", i)
		result := config.parser(request)

		if result.error == nil {
			config.callback(result)
			return result.size, nil
		}

		if result.size > bestResult.size {
			bestResult = result
			bestResultIndex = i
		}
	}

	configs[bestResultIndex].callback(bestResult)
	return bestResult.size, bestResult.error
}
