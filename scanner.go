package main

import (
	"fmt"
	"strconv"
)

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
		s.skipWhiteSpace()
		s.start = s.current
		tok := s.ScanToken()
		tokens = append(tokens, tok)
	}
	tokens = append(tokens, Token{
		Type:   TokenEOF,
		Lexeme: "",
		Line:   s.line,
	})

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

	tkn := Token{
		Type:    TokenString,
		Lexeme:  s.source[s.start:s.current],
		Literal: value,
		Line:    s.line,
	}
	// fmt.Printf("[string] returning: %#v\n", tkn)
	return tkn
}

func (s *Scanner) ScanToken() Token {
	// fmt.Printf("ScanToken start %d: current %d: char %q\n", s.start, s.current, s.peek())
	for {
		if s.isAtEnd() {
			return Token{Type: TokenEOF}
		}
		s.start = s.current
		c := s.advance()

		if isAlpha(c) {
			return s.identifier()
		}

		if isDigit(c) {
			return s.number()
		}

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
		case '\x00':
			continue
		// case ' ':
		// 	return s.ScanToken()
		// case '\n':
		// 	s.line++
		// 	s.advance()
		// 	return s.ScanToken()
		default:
			// if c == '\x00' {
			// 	return Token{Type: TokenEOF, Line: s.line}
			// }
			fmt.Printf("Unhandled char: %q\n", c)
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

func (s *Scanner) skipWhiteSpace() {
	for {
		if s.isAtEnd() {
			return
		}

		switch s.peek() {
		case ' ', '\r', '\t':
			s.advance()
		case '\n':
			s.line++
			s.advance()
		default:
			return
		}
	}
}

func (s *Scanner) identifier() Token {
	fmt.Printf("[identifier] start=%d current=%d\n", s.start, s.current)
	for isAlphaNumeric(s.peek()) {
		s.advance()
		fmt.Printf("[identifier] start=%d current=%d\n", s.start, s.current)
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
	return (c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_')
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
	if s.isAtEnd() {
		return 0
	}
	ch := s.source[s.current]
	s.current++

	return ch
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
