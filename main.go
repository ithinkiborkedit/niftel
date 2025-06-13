package main

import (
	"fmt"
	"log"
	"os"
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
	log.Print("Reading source code...")

	source := string(data)

	lexer := NewScanner(string(source))
	tokens := lexer.ScanTokens()

	parser := NewParser(tokens)
	statements, err := parser.Parse()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Parsing error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Parsed %d statements:\n", len(statements))
	for _, stmt := range statements {
		fmt.Printf("TYPE %T\nDATA: %#v\n", stmt, stmt)
	}

}
