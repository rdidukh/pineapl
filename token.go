package main

import (
	"fmt"
	"regexp"
)

type TokenType int

const (
	TOKEN_TYPE_UNKNOWN TokenType = iota
	TOKEN_TYPE_IDENTIFIER
	TOKEN_TYPE_NUMBER
	TOKEN_TYPE_WHITESPACE
	TOKEN_TYPE_ROUND_BRACKET_OPEN
	TOKEN_TYPE_ROUND_BRACKET_CLOSE
	TOKEN_TYPE_CURLY_BRACKET_OPEN
	TOKEN_TYPE_CURLY_BRACKET_CLOSE
	TOKEN_TYPE_EQUALS
	TOKEN_TYPE_LESS_THAN
	TOKEN_TYPE_COMMA
	TOKEN_TYPE_KEYWORD_FUNC
)

type Token struct {
	value     string
	tokenType TokenType
	start     int
	end       int
}

func (t TokenType) String() string {
	switch t {
	case TOKEN_TYPE_UNKNOWN:
		return "UNKNOWN"
	case TOKEN_TYPE_WHITESPACE:
		return "WHITESPACE"
	case TOKEN_TYPE_IDENTIFIER:
		return "IDENTIFIER"
	case TOKEN_TYPE_NUMBER:
		return "NUMBER"
	case TOKEN_TYPE_ROUND_BRACKET_OPEN:
		return "ROUND_BRACKET_OPEN"
	case TOKEN_TYPE_ROUND_BRACKET_CLOSE:
		return "ROUND_BRACKET_CLOSE"
	case TOKEN_TYPE_CURLY_BRACKET_OPEN:
		return "CURLY_BRACKET_OPEN"
	case TOKEN_TYPE_CURLY_BRACKET_CLOSE:
		return "CURLTY_BRACKET_CLOSE"
	case TOKEN_TYPE_EQUALS:
		return "EQUALS"
	}

	return "UNSUPPORTED"
}

type tokenGetter func(string, int) *Token

func regexTokenGetter(tokenType TokenType, regex *regexp.Regexp) tokenGetter {
	return func(code string, offset int) *Token {
		return getTokenByRegexp(code, offset, tokenType, regex)
	}
}

func runeTokenGetter(tokenType TokenType, rune rune) tokenGetter {
	return func(code string, offset int) *Token {
		return getTokenByRune(code, offset, tokenType, rune)
	}
}

func (t Token) String() string {
	//return fmt.Sprintf("Token(type=%q, start=%d, end=%d, value=%q)", t.tokenType, t.start, t.end, t.value)
	return fmt.Sprintf("%-20s %3d..%2d   %-10q", t.tokenType, t.start, t.end, t.value)
}

func getNextToken(sourceCode string, offset int) *Token {
	tokenGetters := []tokenGetter{
		regexTokenGetter(TOKEN_TYPE_WHITESPACE, regexp.MustCompile(`^\s+`)),
		regexTokenGetter(TOKEN_TYPE_IDENTIFIER, regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*`)),
		// TODO: scientific notation.
		regexTokenGetter(TOKEN_TYPE_NUMBER, regexp.MustCompile(`^[-+]?[0-9]+(\.[0-9]+)?`)),
		runeTokenGetter(TOKEN_TYPE_ROUND_BRACKET_OPEN, '('),
		runeTokenGetter(TOKEN_TYPE_ROUND_BRACKET_CLOSE, ')'),
		runeTokenGetter(TOKEN_TYPE_CURLY_BRACKET_OPEN, '{'),
		runeTokenGetter(TOKEN_TYPE_CURLY_BRACKET_CLOSE, '}'),
		runeTokenGetter(TOKEN_TYPE_EQUALS, '='),
		runeTokenGetter(TOKEN_TYPE_LESS_THAN, '<'),
		runeTokenGetter(TOKEN_TYPE_COMMA, ','),
	}

	for _, tokenGetter := range tokenGetters {
		token := tokenGetter(sourceCode, offset)

		if token.tokenType != TOKEN_TYPE_UNKNOWN {
			return token
		}
	}

	return tokenTypeUnknown
}

func getTokenByRune(sourceCode string, offset int, tokenType TokenType, rune rune) *Token {
	firstRune := getFirstRune(sourceCode)

	if firstRune != rune {
		return tokenTypeUnknown
	}

	return &Token{
		value:     string(rune),
		tokenType: tokenType,
		start:     offset,
		end:       offset + 1, // TODO: can be more than 1.
	}
}

func getTokenByRegexp(sourceCode string, offset int, tokenType TokenType, regex *regexp.Regexp) *Token {
	loc := regex.FindStringIndex(sourceCode)

	if loc == nil {
		return tokenTypeUnknown
	}

	start := loc[0]
	end := loc[1]
	token := sourceCode[start:end]

	if start != 0 {
		panic("start != 0")
	}

	return &Token{
		value:     token,
		tokenType: tokenType,
		start:     offset,
		end:       offset + end,
	}
}

func GetTokens(code string) ([]*Token, error) {
	offset := 0
	tokens := []*Token{}

	for offset < len(code) {
		token := getNextToken(code[offset:], offset)
		if token.tokenType == TOKEN_TYPE_UNKNOWN {
			return tokens, fmt.Errorf("Unknown token at offset %d: %q", offset, code[offset])
		}

		tokens = append(tokens, token)
		offset = token.end
	}

	secondPass(tokens)

	return tokens, nil
}

// TODO: rename.
var keywords = map[string]TokenType{
	"func": TOKEN_TYPE_KEYWORD_FUNC,
}

func secondPass(tokens []*Token) {
	for _, token := range tokens {
		tokenType, ok := keywords[token.value]
		if ok {
			token.tokenType = tokenType
		}
	}
}

func getFirstRune(str string) rune {
	for _, r := range str {
		return r
	}
	panic("empty string")
}
