package main

import (
	"fmt"
	"os"
	"strconv"
)

type TokenType string

const (
	TokenEOF       TokenType = "EOF"
	TokenIdent     TokenType = "IDENT"
	TokenNumber    TokenType = "NUMBER"
	TokenString    TokenType = "STRING"
	TokenTrue      TokenType = "TRUE"
	TokenNil       TokenType = "NIL"
	TokenFalse     TokenType = "FALSE"
	TokenEqal      TokenType = "EQUAL"
	TokenBang      TokenType = "BANG"
	TokenStar      TokenType = "STAR"
	TokenFWDSlash  TokenType = "FWD_SLASH"
	TokenPlus      TokenType = "PLUS"
	TokenMinus     TokenType = "MINUS"
	TokenLParen    TokenType = "LPAREN"
	TokenRParen    TokenType = "RPAREN"
	TokenIf        TokenType = "KEYWORD"
	TokenFor       TokenType = "KEYWORD"
	TokenVar       TokenType = "KEYWORD"
	TokenLBrace    TokenType = "LBRACE"
	TokenRBrace    TokenType = "RBRACE"
	TokenLBracket  TokenType = "LBRACKET"
	TokenRBracket  TokenType = "RBRACKET"
	TokenRepo      TokenType = "REPO"
	TokenBranch    TokenType = "BRANCH"
	TokenIn        TokenType = "KEYWORD"
	TokenElse      TokenType = "KEYWORD"
	TokenEqality   TokenType = "EQUAL_EQUAL"
	TokenBangEqal  TokenType = "BANG_EQUAL"
	TokenGreater   TokenType = "GREATER"
	TokenLess      TokenType = "LESS"
	TokenGreaterEq TokenType = "GREATER_EQUAL"
	TokenLessEq    TokenType = "LESS_EQUAL"
)

var keywords = map[string]TokenType{
	"true":   TokenTrue,
	"nil":    TokenNil,
	"false":  TokenFalse,
	"if":     TokenIf,
	"else":   TokenElse,
	"for":    TokenFor,
	"in":     TokenIn,
	"var":    TokenVar,
	"repo":   TokenRepo,
	"branch": TokenBranch,
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

type Scanner struct {
	source  string
	start   int
	end     int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		start:   0,
		end:     0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() []Token {
	var tokens []Token
	for !s.isAtEnd() {
		s.start = s.current
		tokens = append(tokens, s.ScanToken())
	}
	return tokens
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) string() Token {
	for !s.isAtEnd() && s.peek() != '"' {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return Token{
			Type:   TokenType("Error"),
			Lexeme: "Unterminated string",
			Line:   s.line,
		}
	}
	s.advance()

	value := s.source[s.start+1 : s.current-1]

	return Token{
		Type:    TokenString,
		Lexeme:  s.source[s.start:s.current],
		Literal: value,
		Line:    s.line,
	}
}

func (s *Scanner) ScanToken() Token {
	c := s.advance()
	switch c {
	case '(':
		return s.makeToken(TokenLParen)
	case ')':
		return s.makeToken(TokenRParen)
	case '+':
		return s.makeToken(TokenPlus)
	case '-':
		return s.makeToken(TokenMinus)
	case '=':
		return s.makeToken(TokenEqal)
	case '"':
		return s.string()
	case ' ', '\r', '\t':
		return s.ScanToken()
	case '\n':
		s.line++
		return s.ScanToken()
	default:
		if isDigit(c) {
			return s.number()
		} else if isAlpha(c) {
			return s.identifier()
		} else {
			return s.makeToken(TokenType("ERROR"))
		}
	}
}

func (s *Scanner) number() Token {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	value := s.source[s.start:s.current]
	return Token{
		Type:    TokenNumber,
		Lexeme:  value,
		Literal: parseNumber(value),
		Line:    s.line,
	}
}

func (s *Scanner) identifier() Token {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	if tokenType, ok := keywords[text]; ok {
		return Token{
			Type:   tokenType,
			Lexeme: text,
			Line:   s.line,
		}
	}

	return Token{
		Type:   TokenIdent,
		Lexeme: text,
		Line:   s.line,
	}
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z')
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func parseNumber(text string) float64 {
	val, _ := strconv.ParseFloat(text, 64)
	return val
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) makeToken(t TokenType) Token {
	return Token{
		Type:   t,
		Lexeme: s.source[s.start:s.current],
		Line:   s.line,
	}
}

func main() {

	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: niftel <file.nif>")
		os.Exit(1)
	}

	filename := os.Args[1]

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file %s: %v\n", filename, err)
		os.Exit(1)
	}

	lexer := NewScanner(string(data))
	tokens := lexer.ScanTokens()

	parser := NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parsing error: %v\n", err)
		os.Exit(1)
	}

	for _, stmt := range statements {
		fmt.Printf("%#v\n", stmt)
	}

}
