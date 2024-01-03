package token

import (
	"fmt"
	"regexp"
)

type Type int

const (
	TYPE_UNKNOWN Type = iota
	TYPE_IDENTIFIER
	TYPE_NUMBER
	TYPE_WHITESPACE
	TYPE_ROUND_BRACKET_OPEN
	TYPE_ROUND_BRACKET_CLOSE
	TYPE_CURLY_BRACKET_OPEN
	TYPE_CURLY_BRACKET_CLOSE
	TYPE_EQUALS
	TYPE_LESS_THAN
	TYPE_COMMA
	TYPE_KEYWORD_FUNC
	TYPE_EOF
)

type Token struct {
	Value string
	Type  Type
	Start int
	End   int
}

func (t *Token) IsEof() bool {
	return t.Type == TYPE_EOF
}

func (t Type) String() string {
	switch t {
	case TYPE_UNKNOWN:
		return "UNKNOWN"
	case TYPE_WHITESPACE:
		return "WHITESPACE"
	case TYPE_IDENTIFIER:
		return "IDENTIFIER"
	case TYPE_NUMBER:
		return "NUMBER"
	case TYPE_ROUND_BRACKET_OPEN:
		return "ROUND_BRACKET_OPEN"
	case TYPE_ROUND_BRACKET_CLOSE:
		return "ROUND_BRACKET_CLOSE"
	case TYPE_CURLY_BRACKET_OPEN:
		return "CURLY_BRACKET_OPEN"
	case TYPE_CURLY_BRACKET_CLOSE:
		return "CURLY_BRACKET_CLOSE"
	case TYPE_EQUALS:
		return "EQUALS"
	case TYPE_KEYWORD_FUNC:
		return "KEYWORD_FUNC"
	case TYPE_COMMA:
		return "COMMA"
	case TYPE_LESS_THAN:
		return "LESS_THAN"
	case TYPE_EOF:
		return "EOF"
	}

	panic(fmt.Sprintf("Unsupported token type: %d", int(t)))
}

var tokenTypeUnknown = &Token{
	Type: TYPE_UNKNOWN,
}

type tokenGetter func(string, int) *Token

func regexTokenGetter(tokenType Type, regex *regexp.Regexp) tokenGetter {
	return func(code string, offset int) *Token {
		return getTokenByRegexp(code, offset, tokenType, regex)
	}
}

func runeTokenGetter(tokenType Type, rune rune) tokenGetter {
	return func(code string, offset int) *Token {
		return getTokenByRune(code, offset, tokenType, rune)
	}
}

func getNextToken(sourceCode string, offset int) *Token {
	tokenGetters := []tokenGetter{
		regexTokenGetter(TYPE_WHITESPACE, regexp.MustCompile(`^\s+`)),
		regexTokenGetter(TYPE_IDENTIFIER, regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*`)),
		// TODO: scientific notation.
		regexTokenGetter(TYPE_NUMBER, regexp.MustCompile(`^[-+]?[0-9]+(\.[0-9]+)?`)),
		runeTokenGetter(TYPE_ROUND_BRACKET_OPEN, '('),
		runeTokenGetter(TYPE_ROUND_BRACKET_CLOSE, ')'),
		runeTokenGetter(TYPE_CURLY_BRACKET_OPEN, '{'),
		runeTokenGetter(TYPE_CURLY_BRACKET_CLOSE, '}'),
		runeTokenGetter(TYPE_EQUALS, '='),
		runeTokenGetter(TYPE_LESS_THAN, '<'),
		runeTokenGetter(TYPE_COMMA, ','),
	}

	for _, tokenGetter := range tokenGetters {
		token := tokenGetter(sourceCode, offset)

		if token.Type != TYPE_UNKNOWN {
			return token
		}
	}

	return tokenTypeUnknown
}

func getTokenByRune(sourceCode string, offset int, tokenType Type, rune rune) *Token {
	firstRune := getFirstRune(sourceCode)

	if firstRune != rune {
		return tokenTypeUnknown
	}

	return &Token{
		Value: string(rune),
		Type:  tokenType,
		Start: offset,
		End:   offset + 1, // TODO: can be more than 1.
	}
}

func getTokenByRegexp(sourceCode string, offset int, tokenType Type, regex *regexp.Regexp) *Token {
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
		Value: token,
		Type:  tokenType,
		Start: offset,
		End:   offset + end,
	}
}

func GetTokens(code string) ([]*Token, error) {
	offset := 0
	tokens := []*Token{}

	for offset < len(code) {
		token := getNextToken(code[offset:], offset)
		if token.Type == TYPE_UNKNOWN {
			return tokens, fmt.Errorf("Unknown token at offset %d: %q", offset, code[offset])
		}

		tokens = append(tokens, token)
		offset = token.End
	}

	tokens = append(tokens, &Token{
		Type: TYPE_EOF,
	})
	secondPass(tokens)

	return tokens, nil
}

func GetIterator(code string) (*Iterator, error) {
	tokens, err := GetTokens(code)

	if err != nil {
		return &Iterator{}, err
	}

	return &Iterator{
		tokens: tokens,
	}, nil
}

type Iterator struct {
	tokens []*Token
	index  int
}

func (it Iterator) IsEof() bool {
	return it.index+1 >= len(it.tokens)
}

func (it *Iterator) Advance() {
	it.index += 1
}

func (it *Iterator) Token() *Token {
	if it.index >= len(it.tokens) {
		return &Token{Type: TYPE_EOF}
	}
	return it.tokens[it.index]
}

// TODO: rename.
var keywords = map[string]Type{
	"func": TYPE_KEYWORD_FUNC,
}

func secondPass(tokens []*Token) {
	for _, token := range tokens {
		tokenType, ok := keywords[token.Value]
		if ok {
			token.Type = tokenType
		}
	}
}

func getFirstRune(str string) rune {
	for _, r := range str {
		return r
	}
	panic("empty string")
}
